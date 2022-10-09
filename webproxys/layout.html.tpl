<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/8.0.1/normalize.min.css"
		integrity="sha512-NhSC1YmyruXifcj/KFRWoC561YpHpc5Jtzgvbuzx5VozKpWvQ+4nXhPdFgmx8xqexRcpAglTj9sIBWINXa8x5w=="
		crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="/style.css">
	<title>Gemocities</title>
</head>

<body>
	<nav>
		<h1><a href="/">Gemocities</a></h1>
		{{ if not .Error }}
		<div class="dimmed">
			<p class="notice">
				This is a proxy of Gemini content at <code>{{ .Path }}</code>.
				<br>
				{{ if .UserContent }}
				Gemocities is not responsible for user-generated content.
				<br>
				{{ end }}
				<a href="{{ .GeminiURL }}">View original in Gemini &raquo;</a>
			</p>
		</div>
		{{ end }}
	</nav>
	<main>
		{{ if .Error }}
		<div class="error">
			<p>{{ .Content }}</p>
			<p><a href="{{ .GeminiURL }}">View original in Gemini &raquo;</a></p>
		</div>
		{{ else }}
		{{ .Content }}
		{{ end }}
	</main>
	<footer>
		<div class="dimmed">
			<p>
				<a href="https://github.com/mplewis/gemocities" target="_blank">GitHub</a>
				|
				<a href="mailto:webmaster@gemocities.com">contact</a>
			</p>
		</div>
	</footer>
</body>

</html>
