package templates

const INDEX_TEMPLATE_DISCOVER = `
<html>
<head>
    <title>Welcome To K8Guard Discover - {{.Cluster}}</title>
</head>
<body>
<h2>K8Guard Discover</h2>
<hr/>
<ul>
    {{range .Names}}
    <li><a href="./{{.}}">  {{.}} </a></li>
    {{end}}
</ul>
<hr>
<h6> {{.Version}} </h6>

</body>
</html>`
