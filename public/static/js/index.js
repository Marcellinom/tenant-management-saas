const token = sessionStorage.key("token")
if (token == null) {
    $.ajax({
        // Post username, password & the grant type to /token
        url: "{{.iam_url}}",
        method: 'get',
        contentType: 'application/json',
        //response(Access Token) stores inside session storage of Client browser.
    })
}
else {
    console.log(token)
}