# Flammable Note

Flammable Note is a simple note keeping, stats program built with the following concepts.

## Design Rules

* Control your Data
* Simple as Fuck
* Always Backwards Comptaible
* Today is Blank

Control your data means that everything is stored in a manner you have complete control. At present the design is a single file per day using CSV. No matter what we do with the data a lay person should be able to copy and access the data source with text editor or similiar simple tool.

Simple as Fuck means my non-techy friends should be able to setup and manage their own instace. Also where any decision would add complexity it needs to fuck off or be the bee's knees.
Examples
* What about rich media: Don't care
* What about formatting: Meh
* What about timezomes: Nope
* What about UTF-8: Yeah makes sense to support unicode
* What about version control: Meh

Always backwards compatible. No ifs or butts.

Today is Blank is a key concept to fight depression and task overload. Everything is focused on today by default, with no option to make future notes or persistent entries. This limitation keeps it out of the calender, appointment ect... mindset.

## Future Plan - AI/Grep

This is built stupid simple for a reason because if it fucks up or you don't like it then it's super easy to stop. Though the other reason is I plan to use some smart tools with this from a simple linux Grep to scraping sentiment analysis and smart crap like that. But that is future Claire concern. Right now I need the webserver equivilient of my handbag notebook.

## API v1.0 Stable

### DATA FORMAT
2009-01-30, v1.0
0,My Note,08:30,
1,This note was modified,12:30,13:30
2,,11:00,15:00
3,This is a note with \n a new line\, and a comma,13:00,
4,Notice that line 2 was deleted but we kept the index,10:00
5,Also notice the data at the top and version numer,10:00

### POST /new?timestamp=
* Body text is the new note
* Timestamp is optional if the client wishes to provide an alternate timestamp. It must be in the past but on today's date
* Response is the full CSV for today (text data is cheap)

### POST /edit?index=&timestamp
* Body text is the new note
* Index is a MUST and is the offset for that note
* Timestamp is an optional timestamp if the client wishes to provide. In the case where the timestamp is NOT today you can edit a past note. Though the index must be in said past

### GET /today
* Response is full CSV for today

### GET /past?timestamp=
* Response is the full text from some past day