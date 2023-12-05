# Go Web App Demo

Code for a demo web application written in Go, built as part of Trevor Sawler's excellent online course [Building Modern Web Applications With Go](https://www.udemy.com/course/building-modern-web-applications-with-go)

This web app demo is a fully functional Bed and Breakfast booking website that allows users to:
- View availability of rooms
- Book rooms on an available set of dates
- Get notified about their reservations via email / text

It also allows administrators to:
- View current reservations
- Block out reservations on specified dates
- Get notified about reservations made via email / text

- Built in Go
- Uses the [chi router](https://github.com/go-chi-chi)
- Uses [Alex Edwards SCS](https://github.com/alexedwards/scs/v2) session management
- Uses [nosurf](https://github.com/justinas/nosurf)