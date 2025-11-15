package main

func homeHandler(req Request) Response {
	return Response{
		StatusCode: 200,
		Body:       "Welcome to GHTTP -- HTTP server from raw TCP in Go!",
	}
}

func healthHandler(req Request) Response {
	return Response{
		StatusCode: 200,
		Body:       "OK\n",
	}
}

func echoHandler(req Request) Response {
	body := "Echo: "
	if len(req.Body) > 0 {
		body += string(req.Body)
	} else {
		body += "(No payload provided)"
	}
	body += "\n"
	return Response{
		StatusCode: 200,
		Body:       body,
	}
}
