package importer

var htmlTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <title>Fairlance Importer</title>
        <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
        <style>
            html {
                overflow-y: scroll;
            }
            body {
                padding-top: 20px;                
            }
            .alert-fixed {
                position:fixed; 
                top: 0px;
                left: 30%;
                width: 40%;
                z-index:9999; 
                border-top-left-radius:0px;
                border-top-right-radius:0px;
            }
        </style>
    </head>
    <body>
        <div id="app" class="container">
            <div v-if="msg !== ''" v-bind:class="{ 'alert-success': msg == 'ok', 'alert-warning': msg == 'Running...', 'alert-danger': (msg != 'ok' && msg != 'Running...') }" class="alert alert-fixed">${msg}</div>
            <nav class="navbar navbar-default">
                <div class="container-fluid">
                    <div class="navbar-header">
                        <a class="navbar-brand" href="#">
                            Fairlance Importer
                        </a>
                    </div>
                    <ul class="nav nav-pills">
                        <li v-bind:class="{ 'active': tab == 'db' }">
                            <a href="#" class="navbar-btn" v-on:click="tab = 'db'">DB</a>
                        </li>
                        <li v-bind:class="{ 'active': tab == 'search' }">
                            <a href="#" class="navbar-btn" v-on:click="tab = 'search'">Search</a>
                        </li>
                    </ul>
                </div>
            </nav>
            <div class="row">
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
                        <button type="button" v-bind:class="{ 'btn-primary': type == 'jobs'}" class="btn btn-default" v-on:click="type = 'jobs'">Jobs</button>
                        <button type="button" v-bind:class="{ 'btn-primary': type == 'freelancers'}" class="btn btn-default" v-on:click="type = 'freelancers'">Freelancers</button>
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
                                    <th>updatedAt</th>
                                    <th>action</th>
                                </tr>
                            </thead>
                            <tbody>
                            <template v-for="entity in entities">
                                <tr>
                                    <td>${ entity.id }</td>
                                    <td>${ entity.name }${ entity.firstName } ${ entity.lastName }</td>
                                    <td>${ entity.updatedAt }</td>
                                    <td>
                                        <button type="button" class="btn btn-default btn-sm" v-if="(document !== null && docID != entity.id) || document == null" v-on:click="updateWithDocID('action=get', entity.id)">
                                            <span class="glyphicon glyphicon-chevron-down"></span>
                                        </button>
                                        <!--<button v-if="(document !== null && docID != entity.id) || document == null" v-on:click="updateWithDocID('get', entity.id)">get</button>-->
                                        <button type="button" class="btn btn-default btn-sm" v-if="document !== null && docID == entity.id" v-on:click="document = null">
                                            <span class="glyphicon glyphicon-chevron-up"></span>
                                        </button>
                                        <!--<button v-if="document !== null && docID == entity.id" v-on:click="updateWithDocID('hide', entity.id)">hide</button>-->
                                        <button type="button" class="btn btn-default btn-sm" v-on:click="updateWithDocID('action=import', entity.id)">
                                            <span class="glyphicon glyphicon-upload"></span> import
                                        </button>
                                        <!--<button v-on:click="updateWithDocID('import', entity.id)">import</button>-->
                                        <button type="button" class="btn btn-default btn-sm" v-on:click="updateWithDocID('action=remove', entity.id)">
                                            <span class="glyphicon glyphicon-erase"></span> remove
                                        </button>
                                        <!--<button v-on:click="updateWithDocID('remove', entity.id)">remove</button>-->
                                    </td>
                                </tr>
                                <template v-if="document !== null && docID == entity.id">
                                <tr>
                                    <td colspan="4">
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
                    <div class="well">
                        <div class="btn-group pull-right">
                            <button type="button" v-bind:class="{ 'btn-primary': type == 'jobs'}" class="btn btn-default" v-on:click="type = 'jobs'">Jobs</button>
                            <button type="button" v-bind:class="{ 'btn-primary': type == 'freelancers'}" class="btn btn-default" v-on:click="type = 'freelancers'">Freelancers</button>
                        </div>
                        <h1>${type}</h1>
                        <form class="form-horizontal">
                            <div class="form-group">
                                <div class="col-sm-12 col-md-3">
                                    <input type="number" class="form-control" v-model="docID" placeholder="id">
                                </div>
                            </div>
                            <div class="form-group">
                                <div class="col-sm-11">
                                    <button type="button" class="btn btn-default" v-on:click="search">Search</button>
                                </div>
                            </div>
                        </form>
                    </div>
                    <div v-if="entities !== null" class="row">
                        <div class="col-md-12" v-for="entity in entities">
                            <div class="panel panel-default">
                                <div class="panel-body">
                                    <div v-for="(val, key) in entity.fields">
                                        <span><b>${key}</b>: ${val}</span>
                                    </div>
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
                period: null,
                selectedTag: null,
                tags: [],
                priceFrom: null,
                priceTo: null
            }
        },
        watch: {
                type: function() {
                    if (this.tab == "db") {
                        this.update()
                    }
                },
                tab: function() {
                    this.entities = []
                    if (this.tab == "db") {
                        this.update()
                    }
                }
            },
        methods: {
            updateWithDocID: function(GETParams, id) {this.update(GETParams + "&docID=" + id);},
            update: function(GETParams) {
                this.msg = "Running...";
                var params = 'type=' + this.type + '&tab=' + this.tab;
                if (GETParams != undefined) {
                    if (GETParams[0] != '&') GETParams = "&" + GETParams;
                    params = params + GETParams;
                }
                var app = this;
                axios.get('http://local.fairlance.io:3004/json?' + params)
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
                        } else if (app.tab == 'search') {
                            app.search.tags = response.data.Search.Tags;
                        }
                        window.setTimeout(function() {
                            app.msg = '';
                        }, 1000);
                    })
                    .catch(function (error) {
                        console.log(error)
                        app.msg = "Error"
                    })
            },
            prevPageLabel: function() {
                if (this.offset >= this.limit) {
                    return (this.offset-this.limit+1) + "-" + this.offset
                }

                return ""
            },
            nextPageLabel: function() {
                if (this.offset + this.limit < this.totalInDB) {
                    return (this.offset+this.limit+1) + "-" + (this.offset+(this.limit*2))
                }

                return ""
            },
            currentPageLabel: function() {
                return (this.offset + 1) + "-" + (this.offset+this.limit)
            },
            search: function() {
                var params = {
                    // period: this.search.period,
                    // priceFrom: this.search.priceFrom,
                    // priceTo: this.search.priceTo,
                    // tags: this.search.selectedTag,
                    docID: this.docID
                }
                this.update('action=search&' + this.serialize(params))
            },
            serialize: function(obj) {
                var str = [];
                for(var p in obj)
                    if (obj.hasOwnProperty(p) && obj[p] != undefined) {
                    str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
                    }
                return str.join("&");
            }
        }
    })
    app.update();
  </script>
</html>
`
