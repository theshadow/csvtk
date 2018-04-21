# The CSV Toolkit (csvtk)

This tool aims to provide similar functionality to csvkit. Why duplicate this project? Because I didn't like having
to deploy the entire python binary and dependencies everywhere I wanted to use the tool. So, here we are.

To start CSV Toolkit comes with rendering functionality. You can pipe in CSV data via Stdin and have it render a
nicely formatted ASCII table.

# Documentation

## render

Render allows you to render CSV data in a human readable ASCII formatted table. By default it will consume CSV
formatted records from Stdin. You may specify a file by using the **--input** flag. Similarly, output will write
to Stdout by default and may be changed with the **--output** flag

### Examples

- Default :: `cat MOCK_DATA.csv | csvtk render`
- With input :: `csvtk render --input MOCK_DATA.csv`
- With output :: `cat MOCK_DATA.csv | csvtk render --output MOCK_DATA.txt`

## Select

Select allows you to specify which columns you'd like to display. You do this by

- IF(), LEN(), UPPER(), LOWER(), ($1+$2), IFNULL(), GT(), LT(), EQ(), GTE(), LTE(), LEFT(), RIGHT(), SUBSTR(),
  AND(), OR(), NOT(), XOR(), BIT_RSHIFT(), BIT_LSHIFT()
- SUM(), AVG(), STD()

"
SELECT $1 AS fname, IF(GT($2, $3), $2, $3)
WHERE (fname > 5)
GROUP BY ($1)
"

- Specifying columns :: `csvtk select $1,$2,$5`
- Specifying using ranges :: `csvtk select $1:$6`
- Specifying using negative ranges :: `csvtk select -$1`
`csvtk select IF($1, 'False')'`


=SUM($1:$5)
=$1+$3
=$1+$3;'###.##'