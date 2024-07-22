package dax

var templateEnd = `%sVAR link = CONCATENATEX(FILTER({%s}, [Value] <> BLANK()), [Value], "%%20and%%20")
VAR fullLink = IF(link = BLANK(), BLANK(), [Dashboard Link] & "?rs:embed=true&filter=" & link)
RETURN
	fullLink`

var templateString = `VAR V%d =
    IF(
        ISFILTERED( %s ) = TRUE(),
        VAR vals = SELECTCOLUMNS(ADDCOLUMNS(VALUES( %s ), "clean",
            SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(%s,
                "%%", "%%25"),
                " ", "%%20"),
                "'", "%%27%%27"),
                "+", "%%2B"),
                "/", "%%2F"),
                "?", "%%3F"),
                "#", "%%23"),
                "&", "%%26")
            ), "clean", [clean])
        VAR cnt = COUNTROWS(vals)
        VAR valsStr = IF(cnt = 1, "eq%%20%%27" & vals & "%%27", "in%%20(%%27" & CONCATENATEX(vals, [clean], "%%27,%%20%%27") & "%%27)")
        VAR fullStr = "%s/%s%%20" & valsStr
        VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
    	RETURN
        	fullStrBlank,   
        	BLANK()
    )

`

var templateDate = `VAR V%d =
    IF(
		ISFILTERED( %s ) = TRUE(),
		VAR vals = VALUES( %s )
		VAR cnt = COUNTROWS(vals)
        VAR minDate = FIRSTDATE(vals)
        VAR maxDate = LASTDATE(vals)
        var daysBetweenMinMax = DATEDIFF(minDate, maxDate, DAY)
        VAR valsStr = IF(cnt = 1, "eq%%20datetime%%27" & FORMAT(vals, "YYYY-MM-DD") & "T00:00:00%%27",
            IF(daysBetweenMinMax = cnt - 1,
                "ge%%20datetime%%27" & FORMAT(minDate, "YYYY-MM-DD") & "T00:00:00%%27%%20and%%20%s/%s%%20le%%20datetime%%27" & FORMAT(maxDate, "YYYY-MM-DD")  & "T00:00:00%%27",
                "in%%20(datetime%%27" & CONCATENATEX(vals, FORMAT(%s, "YYYY-MM-DD") & "T00:00:00", "%%27,%%20datetime%%27") & "%%27)"))
        VAR fullStr = "%s/%s%%20" & valsStr
        VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
    	RETURN
        	fullStrBlank,
        	BLANK()
    )

`

var templateInt = `VAR V%d =
	IF(
		ISFILTERED( %s ) = TRUE(),
		VAR vals = VALUES( %s )
		VAR cnt = COUNTROWS(vals)
		VAR minInt = MIN(%s)
		VAR maxInt = MAX(%s)
		var betweenMinMax = maxInt - minInt
		VAR valsStr = IF(cnt = 1, "eq%%20" & minInt,
		IF(betweenMinMax = cnt - 1,
			"ge%%20" & minInt & "%%20and%%20%s/%s%%20le%%20" & minInt,
			"in%%20(" & CONCATENATEX(vals, %s, ",%%20") & ")"))
		VAR fullStr = "%s/%s%%20" & valsStr
		VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
		RETURN
			fullStrBlank,
			BLANK()
	)

`
