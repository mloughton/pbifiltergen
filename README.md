# pbifiltergen

A tool that generates a DAX measure for creating pre-filtered Power BI dashboard links.

## Examples

Input:
```
TABLE[COLUMN],STRING
'TABLE 2'[COLUMN 2],DATE
```

Output:
```dax
VAR V0 =
    IF(
        ISFILTERED( TABLE[COLUMN] ) = TRUE(),
        VAR vals = SELECTCOLUMNS(ADDCOLUMNS(VALUES( TABLE[COLUMN] ), "clean",
            SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(TABLE[COLUMN],
                "%", "%25"),
                " ", "%20"),
                "'", "%27%27"),
                "+", "%2B"),
                "/", "%2F"),
                "?", "%3F"),
                "#", "%23"),
                "&amp;", "%26")
            ), "clean", [clean])
        VAR cnt = COUNTROWS(vals)
        VAR valsStr = IF(cnt = 1, "eq%20%27" &amp; vals &amp; "%27", "in%20(%27" &amp; CONCATENATEX(vals, [clean], "%27,%20%27") &amp; "%27)")
        VAR fullStr = "TABLE/COLUMN%20" &amp; valsStr
        VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
    	RETURN
        	fullStrBlank,   
        	BLANK()
    )

VAR V1 =
    IF(
		ISFILTERED( 'TABLE 2'[COLUMN 2] ) = TRUE(),
		VAR vals = VALUES( 'TABLE 2'[COLUMN 2] )
		VAR cnt = COUNTROWS(vals)
        VAR minDate = FIRSTDATE(vals)
        VAR maxDate = LASTDATE(vals)
        var daysBetweenMinMax = DATEDIFF(minDate, maxDate, DAY)
        VAR valsStr = IF(cnt = 1, "eq%20datetime%27" &amp; FORMAT(vals, "YYYY-MM-DD") &amp; "T00:00:00%27",
            IF(daysBetweenMinMax = cnt - 1,
                "ge%20datetime%27" &amp; FORMAT(minDate, "YYYY-MM-DD") &amp; "T00:00:00%27%20and%20TABLE_x0020_2/COLUMN_x0020_2%20le%20datetime%27" &amp; FORMAT(maxDate, "YYYY-MM-DD")  &amp; "T00:00:00%27",
                "in%20(datetime%27" &amp; CONCATENATEX(vals, FORMAT('TABLE 2'[COLUMN 2], "YYYY-MM-DD") &amp; "T00:00:00", "%27,%20datetime%27") &amp; "%27)"))
        VAR fullStr = "TABLE_x0020_2/COLUMN_x0020_2%20" &amp; valsStr
        VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
    	RETURN
        	fullStrBlank,
        	BLANK()
    )

VAR link = CONCATENATEX(FILTER({V0,V1}, [Value] &lt;&gt; BLANK()), [Value], "%20and%20")
VAR fullLink = IF(link = BLANK(), BLANK(), [Dashboard Link] &amp; "?rs:embed=true&amp;filter=" &amp; link)
RETURN
	fullLink
```

## Implementation

* HTTP server and routes built using [Go http standard library](https://pkg.go.dev/net/http)

* DAX generation custom built in Go

* Front end built using [HTMX](https://htmx.org/)
