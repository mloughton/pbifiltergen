# pbifiltergen

A tool that generates a DAX measure for creating pre-filtered Power BI dashboard links.

## Examples

Input:
```
TABLE[COLUMN],STRING
```

Output:
```dax
VAR V0 =
    IF(
        ISFILTERED( TABLE[COLUMN] ) = TRUE(),
        VAR vals = SELECTCOLUMNS(ADDCOLUMNS(VALUES( TABLE[COLUMN] ), "clean",
            SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(TABLE[COLUMN],
                " ", "%20"),
                "'", """"),
                "%", "%25"),
                "+", "%2B"),
                "/", "%2F"),
                "?", "%3F"),
                "#", "%23"),
                "&", "%26")
            ), "clean", [clean])
        VAR cnt = COUNTROWS(vals)
        VAR valsStr = IF(cnt = 1, "eq%20%27" & vals & "%27", "in%20(%27" & CONCATENATEX(vals, [clean], "%27,%20%27") & "%27)")
        VAR fullStr = "TABLE/COLUMN%20" & valsStr
        VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
    	RETURN
        	fullStrBlank,   
        	BLANK()
    )

VAR link = CONCATENATEX(FILTER({V0}, [Value] <> BLANK()), [Value], "%20and%20")
VAR fullLink = IF(link = BLANK(), BLANK(), [Dashboard Link] & "?rs:embed=true&filter=" & link)
RETURN
	fullLink
```

## Implementation

* HTTP server and routes built using [Go http standard library](https://pkg.go.dev/net/http)

* DAX generation custom built in Go

* Front end built using [HTMX](https://htmx.org/)