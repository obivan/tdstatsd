package main

import "net/http"

const indexPage = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="/static/bootstrap.css">
		<script src="/static/vue.js"></script>
		<title>Pools status</title>
	</head>
	<body>
		<div class="container">
			<div class="row">
				<div class="col">
					<div id="app">
						<template v-if="error">
							<div class="alert alert-danger" role="alert">{{ error }}</div>
						</template>
						<template v-else>
							<div class="shadow p-3 mb-5 rounded">
								<span class="badge badge-pill badge-primary">Last updated at {{ new Date().toLocaleString() }}</span>
								<span class="badge badge-pill badge-primary">Count of pools: {{ this.pools !== null ? this.pools.length : 0 }}</span>
								<span class="badge badge-pill badge-primary">Poll interval: {{ this.pollInterval }}s</span>
							</div>
							<table class="table table-borderless table-striped shadow rounded">
								<thead>
									<tr>
										<th scope="col">Name</th>
										<th scope="col">URL</th>
										<th scope="col">Status</th>
									</tr>
								</thead>
								<tbody>
								<tr v-for="pool in pools">
									<td>{{ pool.Name }}</td>
									<td>{{ pool.URL }}</td>
									<template v-if="pool.Status === 'online'">
										<td class="table-success rounded">{{ pool.Status }}</td>
									</template>
									<template v-else>
										<td class="table-danger rounded">{{ pool.Status }}</td>
									</template>
								</tr>
								</tbody>
							</table>
						</template>
					</div>
				</div>
			</div>
		</div>
		<script>
			var app = new Vue({
			el: '#app',
				data: {
					pools: null,
					error: null,
					pollInterval: 5
				},
				methods: {
					loadData: function () {
						fetch("/pools")
							.then(res => {
								if (!res.ok) { throw res }
								return res.json()
							})
							.then(json => {
								this.error = null;
								this.pools = json;
							})
							.catch(err => {
								if (typeof err.text === 'function') {
									err.text().then(errMsg => {
										console.error(errMsg);
										this.error = 'tdstatsd server error: ' + errMsg;
										this.pools = null;
									});
								} else {
									this.error = 'Something went wrong. Make sure the server is running.';
									this.pools = null;
								}
							});
					},
					setUpdateInterval: function() {
						setInterval(function() {
							this.loadData();
						}.bind(this), this.pollInterval * 1000);
					}
				},
				created: function() {
					this.loadData();
					this.setUpdateInterval();
				}
			});
		</script>
	</body>
</html>
`

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound),
			http.StatusNotFound)
		return
	}
	if _, err := w.Write([]byte(indexPage)); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
}
