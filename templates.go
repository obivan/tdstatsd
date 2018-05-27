package main

const indexPageTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Pools status</title>
		<style type="text/css">
			table {
				font-size: 22pt;
				border-collapse: collapse;
				width: 100%;
			}
			td, th {
				padding: 5px;
				border: thin solid black;
			}
		</style>
	</head>
	<body>
		<div>
			<table>
				<tr>
					<th>Name</th>
					<th>URL</th>
					<th>Status</th>
				</tr>
				{{ range . }}
				<tr>
					<td>{{ .Name }}</td>
					<td>{{ .URL }}</td>
					{{ if ne .Status "online" }}
					<td style="color:red;">{{ .Status }}</td>
					{{ else }}
					<td style="color:green;">{{ .Status }}</td>
					{{ end }}
				{{ end }}
				</tr>
			</table>
		</div>
	</body>
</html>
`
