package importer

var htmlTemplate = `
<!DOCTYPE html>
<html>

<head>
    <title>Fairlance Importer</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u"
        crossorigin="anonymous">
    <style>
            html {
                overflow-y: scroll;
            }
            /*body {
                padding-top: 20px;                
            }*/
            .alert-fixed {
                position:fixed; 
                bottom: 50px;
                right: 30px;
                width: 30%;
                z-index:9999; 
                border-top-left-radius:0px;
                border-top-right-radius:0px;
            }
        </style>
</head>

<body>
    <div id="app">
        <div v-if="msg !== ''" v-bind:class="{ 'alert-success': msg == 'ok', 'alert-warning': msg == 'Running...', 'alert-danger': (msg != 'ok' && msg != 'Running...') }"
            class="alert alert-fixed">${msg}</div>
        <nav class="navbar navbar-default navbar-static-top">
            <div class="container">
                <div class="navbar-header">
                    <a class="navbar-brand" href="#">
                            Fairlance Importer
                        </a>
                </div>
                <ul class="nav navbar-nav">
                    <li v-bind:class="{ 'active': tab == 'db' }">
                        <a href="#" v-on:click="tab = 'db'">DB</a>
                    </li>
                    <li v-bind:class="{ 'active': tab == 'search' }">
                        <a href="#" v-on:click="tab = 'search'">Search</a>
                    </li>
                </ul>
            </div>
        </nav>
        <div class="container">
            <div class="rpxow">
                <div class="col-md-3" v-if="tab == 'db'">
                    <div class="btn-group-vertical btn-block">
                        <button class="btn btn-default btn-sm" v-on:click="update('action=re_generate_test_data')" type="button">Generate test data</button>
                        <button class="btn btn-default btn-sm" v-on:click="update('action=import_all')" type="button">Import All</button>
                        <button class="btn btn-default btn-sm" v-on:click="update('action=delete_all_from_search_engine')" type="button">Delete All From Search Engine</button>
                        <button class="btn btn-default btn-sm btn-danger" v-on:click="update('action=delete_all_from_db')" type="button">Delete all from DB</button>
                    </div>
                </div>
                <div class="col-md-9" v-if="tab == 'db'">
                    <div class="btn-group pull-right">
                        <button type="button" v-bind:class="{ 'btn-warning': type == 'jobs'}" class="btn btn-default" v-on:click="type = 'jobs'">Jobs</button>
                        <button type="button" v-bind:class="{ 'btn-warning': type == 'freelancers'}" class="btn btn-default" v-on:click="type = 'freelancers'">Freelancers</button>
                    </div>
                    <h1>${type}</h1>
                    <div class="clearfix">
                        <span>Total In DB: <span class="badge">${totalInDB}</span></span>
                        <span>TotalIn Search Engine: <span class="badge">${totalInSearchEngine}</span></span>
                        <span class="pull-right">
                                <template v-if="prevPageLabel() !== ''">
                                    <button class="btn btn-default btn-xs" v-on:click="update('offset='+ (offset - limit) + '&limit=' + limit)">${prevPageLabel()}</button>
                                </template>
                                ${currentPageLabel()}
                                <template v-if="nextPageLabel() !== ''">
                                    <button class="btn btn-default btn-xs" v-on:click="update('offset='+ (offset + limit) + '&limit=' + limit)">${nextPageLabel()}</button>
                                </template>
                            </span>
                    </div>
                    <div class="table-responsive">
                        <table class="table table-striped table-condensed">
                            <thead>
                                <tr>
                                    <th>id</th>
                                    <th>name</th>
                                    <th>info</th>
                                    <th>updatedAt</th>
                                    <th>action</th>
                                </tr>
                            </thead>
                            <tbody>
                                <template v-for="entity in entities">
                                    <tr>
                                        <td>${ entity.id }</td>
                                        <td>${ entity.name }${ entity.firstName } ${ entity.lastName }</td>
                                        <td>${ entity.summary }${ entity.email }</td>
                                        <td v-bind:title="entity.updatedAt">${ timeSince(entity.updatedAt) } ago</td>
                                        <td>
                                            <button type="button" class="btn btn-default btn-sm" v-if="(document !== null && docID != entity.id) || document == null"
                                                v-on:click="update('action=get&docID=' + entity.id + '&offset='+ offset + '&limit=' + limit)">
                                                <span class="glyphicon glyphicon-chevron-down"></span>
                                            </button>
                                            <button type="button" class="btn btn-default btn-sm" v-if="document !== null && docID == entity.id" v-on:click="document = null">
                                                <span class="glyphicon glyphicon-chevron-up"></span>
                                            </button>
                                            <button type="button" class="btn btn-default btn-sm" v-on:click="update('action=import&docID=' + entity.id)">
                                                <span class="glyphicon glyphicon-upload"></span> import
                                            </button>
                                            <button type="button" class="btn btn-default btn-sm" v-on:click="update('action=remove&docID=' + entity.id)">
                                                <span class="glyphicon glyphicon-erase"></span> remove
                                            </button>
                                        </td>
                                    </tr>
                                    <template v-if="document !== null && docID == entity.id">
                                        <tr>
                                            <td colspan="5">
                                                <table class="u-full-width">
                                                    <thead></thead>
                                                    <tbody>
                                                        <tr v-for="(val, key) in document">
                                                            <td valign="top"><b>${key}</b></td>
                                                            <td valign="top">${val}</td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </td>
                                        </tr>
                                    </template>
                                </template>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="col-sm-12" v-if="tab == 'search'">
                    <div class="row">
                        <div class="col-sm-12">
                            <div class="well">
                                <div class="row">
                                    <form class="form-horizontal">
                                        <div class="col-md-6">
                                            <div class="form-group">
                                                <div class="col-sm-12">
                                                    <div class="btn-group">
                                                        <button type="button" v-bind:class="{ 'btn-warning': search.method == 'get'}" class="btn btn-default" v-on:click="search.method = 'get'">GET</button>
                                                        <button type="button" v-bind:class="{ 'btn-warning': search.method == 'post'}" class="btn btn-default" v-on:click="search.method = 'post'">POST</button>
                                                        <button type="button" v-bind:class="{ 'btn-warning': search.method == 'put'}" class="btn btn-default" v-on:click="search.method = 'put'">PUT</button>
                                                        <button type="button" v-bind:class="{ 'btn-warning': search.method == 'delete'}" class="btn btn-default" v-on:click="search.method = 'delete'">DELETE</button>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="form-group">
                                                <div class="col-sm-12">
                                                    <div class="input-group">
                                                        <span class="input-group-addon">/api/</span>
                                                        <input type="text" class="form-control" v-on:keyup.enter="doSearch" v-model="search.url" placeholder="url" />
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="col-md-6" v-if="search.method == 'post'">
                                            <div class="form-group">
                                                <div class="col-sm-12">
                                                    <textarea class="form-control" rows="7" v-model="search.body" placeholder="body"></textarea>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="col-sm-12">
                                            <div class="form-group">
                                                <div class="col-sm-12">
                                                    <button type="button" class="btn btn-default" v-on:click="doSearch">Execute</button>
                                                </div>
                                            </div>
                                        </div>
                                    </form>
                                </div>
                            </div>
                            <div class="well" v-if="search.rawData != ''">
                                <pre>${search.rawData}</pre>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
<script src="https://unpkg.com/vue@2.0.3/dist/vue.min.js"></script>
<script src="https://unpkg.com/axios@0.12.0/dist/axios.min.js"></script>

<script>
    var app = new Vue({
        delimiters: ['${', '}'],
        el: '#app',
        data: {
            tab: 'db',
            type: 'jobs',
            msg: '',
            entities: [],
            totalInSearchEngine: 0,
            totalInDB: 0,
            docID: 0,
            document: {},
            offset: 0,
            limit: 0,
            search: {
                body: '',
                method: 'get',
                url: '',
                rawData: ''
            }
        },
        watch: {
            type: function () {
                if (this.tab == "db") {
                    this.update()
                }
            },
            tab: function () {
                this.entities = []
                if (this.tab == "db") {
                    this.update()
                }
            }
        },
        methods: {
            update: function (GETParams) {
                this.msg = "Running...";
                var params = 'type=' + this.type + '&tab=' + this.tab;
                if (GETParams != undefined) {
                    if (GETParams[0] != '&') GETParams = "&" + GETParams;
                    params = params + GETParams;
                }
                var app = this;
                axios.get(location.origin + '/json?' + params)
                    .then(function (response) {
                        app.msg = response.data.Message;
                        app.entities = response.data.Entities;
                        app.tab = response.data.Tab;
                        app.type = response.data.Type;
                        app.offset = response.data.Offset;
                        app.limit = response.data.Limit;
                        app.docID = response.data.DocID;
                        if (app.tab == 'db') {
                            app.totalInSearchEngine = response.data.DB.TotalInSearchEngine;
                            app.totalInDB = response.data.DB.TotalInDB;
                            app.document = response.data.DB.Document;
                        }
                        app.clearMsg();
                    })
                    .catch(function (error) {
                        console.log(error)
                        app.msg = "Error";
                        app.clearMsg();
                    })
            },
            prevPageLabel: function () {
                if (this.offset >= this.limit) {
                    return (this.offset - this.limit + 1) + "-" + this.offset
                }

                return ""
            },
            nextPageLabel: function () {
                if (this.offset + this.limit < this.totalInDB) {
                    return (this.offset + this.limit + 1) + "-" + (this.offset + (this.limit * 2))
                }

                return ""
            },
            currentPageLabel: function () {
                return (this.offset + 1) + "-" + (this.offset + this.limit)
            },
            doSearch: function () {
                this.msg = "Running...";
                var params = {
                    type: this.type,
                    tab: this.tab,
                    action: 'search',
                    url: this.search.url,
                    method: this.search.method,
                    body: this.search.body,
                }
                axios.post(location.origin + '/json', params)
                    .then(function (response) {
                        app.msg = response.data.Message;
                        app.search.rawData = JSON.stringify(response.data.RawData, null, 4);
                        app.clearMsg();
                    })
                    .catch(function (error) {
                        console.log(error)
                        app.msg = "Error";
                        app.clearMsg();
                    })
            },
            clearMsg: function() {
                var app = this;
                window.setTimeout(function () {
                    app.msg = '';
                }, 1000);
            },
            serialize: function (obj) {
                var str = [];
                for (var p in obj)
                    if (obj.hasOwnProperty(p) && obj[p] != undefined) {
                        str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
                    }
                return str.join("&");
            },
            timeSince: function (date) {
                date = new Date(date)
                var seconds = Math.floor((new Date() - date) / 1000);

                var interval = Math.floor(seconds / 31536000);

                if (interval > 1) {
                    return interval + " years";
                }
                interval = Math.floor(seconds / 2592000);
                if (interval > 1) {
                    return interval + " months";
                }
                interval = Math.floor(seconds / 86400);
                if (interval > 1) {
                    return interval + " days";
                }
                interval = Math.floor(seconds / 3600);
                if (interval > 1) {
                    return interval + " hours";
                }
                interval = Math.floor(seconds / 60);
                if (interval > 1) {
                    return interval + " minutes";
                }
                return Math.floor(seconds) + " seconds";
            }
        }
    })
    app.update();
</script>

</html>
`
