package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CmdResult result of running a textual command - split into lines
type CmdResult struct {
	Lines []string
}

// NewCommand run a command and parse the output into lines (CmdResult)
func NewCommand(cmd []string) (*CmdResult, error) {
	clicmd := exec.Command(cmd[0], cmd[1:]...)
	output, err := clicmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	var str = string(output)

	//strange: this one gives an empty line at the end. Don't know what string.Split is doing
	//var ret = strings.Split(str, "\n")

	// FieldsFunc solves the issue - but you need to do a callback that tells it when a character is a delimiter
	f := func(c rune) bool { return c == '\n' }
	var ret = strings.FieldsFunc(str, f)
	return &CmdResult{ret}, err
}

func testCommand() {
	cmd := []string{"ls", "/etc/"}

	lines, err := NewCommand(cmd)
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("no-error:")

		for pos, line := range lines.Lines {
			fmt.Printf("pos: %d line \"%s\"\n", pos, line)
		}
	}
}

func testCommand2() {
	cmd := []string{"kubectl", "explain", "--recursive=true", "pod"}

	lines, err := NewCommand(cmd)
	if err != nil {
		fmt.Println("error: ", err, "cmd:", cmd, "lines:", lines)
	} else {
		fmt.Println("no-error:")

		for pos, line := range lines.Lines {
			fmt.Printf("pos: %d line \"%s\"\n", pos, line)
		}
	}
}

// NewCommandWithTimeout Run a command - but if it does not return within specifiec number of seconds then return a timeout
func NewCommandWithTimeout(cmd []string, timeoutInSeconds int) (*CmdResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000*time.Duration(timeoutInSeconds))
	defer cancel()

	clicmd := exec.CommandContext(ctx, cmd[0], cmd[1:]...)

	output, err := clicmd.CombinedOutput()
	if output != nil {
		var str = string(output)

		f := func(c rune) bool { return c == '\n' }
		var ret = strings.FieldsFunc(str, f)
		return &CmdResult{ret}, err
	}
	return nil, err

}

func testCommandWithTimeout() {
	cmd := []string{"sleep", "100"}

	res, err := NewCommandWithTimeout(cmd, 3)
	if err != nil {
		fmt.Println("first cmd: error: ", err)
	}

	fmt.Println("start second command")

	cmdd := []string{"ls", "/etc"}

	res, err = NewCommandWithTimeout(cmdd, 3)
	if err != nil {
		fmt.Println("second cmd: error: ", err)
	}
	if res != nil {
		fmt.Println("second command output:")
		for pos, line := range res.Lines {
			fmt.Println("lino: ", pos, " line: ", line)
		}
	}

	fmt.Println("eof eof")

}

// CmdColumn a column of the tabular output
type CmdColumn struct {
	ColumnName string
	StartPos   int
	EndPos     int
}

func (c CmdColumn) String() string {
	return fmt.Sprintf("Column{column: %s start: %d end: %d}", c.ColumnName, c.StartPos, c.EndPos)
}

// CmdRow the fields of a row of tabular output
type CmdRow []string

// CmdTable parsed output of a tabular command
type CmdTable struct {
	Columns []CmdColumn
	Rows    []CmdRow
}

func parseColumns(columnDef string) ([]CmdColumn, error) {

	var columnNames = strings.Fields(columnDef)

	var retDef = make([]CmdColumn, len(columnNames))

	for pos, element := range columnNames {
		retDef[pos].ColumnName = element
		retDef[pos].StartPos = strings.Index(columnDef, element)
		if pos > 0 {
			retDef[pos-1].EndPos = retDef[pos].StartPos - 1
		}
	}
	retDef[len(columnNames)-1].EndPos = -1 //len(columnDef)

	//fmt.Println(retDef)

	return retDef, nil
}

func parseRows(rows []string, columns []CmdColumn) ([]CmdRow, error) {

	var rval = make([]CmdRow, len(rows))

	for rowPos, row := range rows {
		var rowFields = make([]string, len(columns))

		for columnPos, columnDef := range columns {

			var columnVal string

			if columnDef.EndPos != -1 {
				columnVal = row[columnDef.StartPos:columnDef.EndPos]
			} else {
				columnVal = row[columnDef.StartPos:]
			}

			columnVal = strings.Trim(columnVal, " ")

			rowFields[columnPos] = columnVal
		}

		rval[rowPos] = rowFields
	}

	return rval, nil
}

// NewTable parse command output into a table (for commands like kubectl )
func NewTable(cmd []string) (*CmdTable, error) {
	cmdRes, err := NewCommand(cmd)

	if err != nil {
		return nil, err
	}

	columnNames, err := parseColumns(cmdRes.Lines[0])

	rows, err := parseRows(cmdRes.Lines[1:], columnNames)
	return &CmdTable{columnNames, rows}, nil
}

func escapeLine(in string) string {

	rval := strings.Replace(in, "<", "&lt;", -1)

	rval = strings.Replace(rval, ">", "&gt", -1)

	return rval
}

func showTable() {

	f, err := os.Create("out.html")
	if err != nil {
		fmt.Print("Can't create file", err)
		return
	}

	styleDef := `<style>
	table {
	  width: 100%;
	  background-color: #FFFFFF;
	  border-collapse: collapse;
	  border-width: 2px;
	  border-color: #7ea8f8;
	  border-style: solid;
	  color: #000000;
	}
	
	table td, table th {
	  border-width: 2px;
	  border-color: #7ea8f8;
	  border-style: solid;
	  padding: 5px;
	}
	
	table thead {
	  background-color: #7ea8f8;
	}
	</style>`

	fmt.Fprint(f, styleDef)

	fmt.Fprint(f, "<h1>K8s api-resources and explanations</h1>")
	versionstr, errr := NewCommand([]string{"kubectl", "version", "--short"})
        if errr != nil {
                fmt.Printf("Can't run version", errr)
                return
        }
	fmt.Fprintln(f, "<pre>")
	for _, verstr := range versionstr.Lines {
		fmt.Fprintln(f, escapeLine(verstr))
	}
	fmt.Fprintln(f, "</pre>")

	var args = []string{"kubectl", "api-resources", "-o", "wide"}

	table, err := NewTable(args)
	if err != nil {
		fmt.Printf("sorry, error %s", err)
		return
	}

	fmt.Fprint(f, "<table border='1'><thead><tr>")
	for _, column := range table.Columns {
		fmt.Fprint(f, "<th>", column.ColumnName, "</th>")
	}
	fmt.Fprint(f, "</tr></thead><tbody>")

	for rowPos, rowData := range table.Rows {

		cmd := []string{"kubectl", "explain", rowData[0]}
		_, err := NewCommand(cmd)
		hasDetails := err == nil

		for colPos, col := range rowData {

			//if hasDetails && (colPos == 0 || colPos == 4) {
			if hasDetails && (table.Columns[colPos].ColumnName == "NAME" || table.Columns[colPos].ColumnName == "KIND") {
				fmt.Fprintf(f, "<td><a href=\"#%d\">%s</a></td>", rowPos, col)
			} else {
				fmt.Fprint(f, "<td>", col, "</td>")
			}
		}
		fmt.Fprint(f, "</tr>\n")
	}
	fmt.Fprint(f, "</tbody></table>")

	for posRow, rowData := range table.Rows {

		entityName := rowData[0]

		cmd := []string{"kubectl", "explain", entityName}
		lines, err := NewCommand(cmd)
		if err == nil {
			fmt.Fprintf(f, "<a name=\"%d\"></a><h3>%s</h3>", posRow, entityName)

			fmt.Fprintf(f, "<table border='1'><thead><tr><th>%s</th></tr></thead><tbody><tr><td><pre>", entityName)

			for _, line := range lines.Lines {
				fmt.Fprintf(f, "%s\n", escapeLine(line))

			}
			fmt.Fprint(f, "</pre></td></tr>")

			cmdd := []string{"kubectl", "explain", "--recursive=true", entityName}
			lines, err = NewCommandWithTimeout(cmdd, 5)
			if err == nil {

				fmt.Fprintf(f, "<tr><th>detailed %s</th></tr><tr><td><pre>", entityName)

				for _, line := range lines.Lines {
					fmt.Fprintf(f, "%s\n", escapeLine(line))
				}

				fmt.Fprint(f, "</pre></td></tr>")
			} else {
				fmt.Println("Command: ", cmdd, " Failed: ", err)
			}
			fmt.Fprint(f, "</tbody></table>")
		} else {
			fmt.Println("Command: ", cmd, " Failed: ", err)
		}
	}
	f.Close()
}

func main() {
	//testCommandWithTimeout()
	//testCommand()
	showTable()
}
