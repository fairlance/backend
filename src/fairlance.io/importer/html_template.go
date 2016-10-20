package importer

var htmlTemplate = `
<html>
  <head>
    <title>Fairlance Importer</title>
    <link href="http://getskeleton.com/dist/css/normalize.css" rel="stylesheet" type="text/css"/>
    <link href="http://getskeleton.com/dist/css/skeleton.css" rel="stylesheet" type="text/css"/>
  </head>
  <body>
    <div class="container">
      <div class="row">
        <div class="columns eleven">
          <b>Status: {{.Message}}</b>
        </div>
      </div>
      <div class="row">
        <div class="column one"><a href="/">Reload</a></div>
        <div class="columns three"><a href="?action=re_generate_test_data&type={{.Type}}">Re-generate test data</a></div>
        <div class="columns three"><a href="?action=import_all&type={{.Type}}">Import to search engine</a></div>
        <div class="columns two"><a href="?action=delete_all_from_db&type={{.Type}}">Delete all from DB</a></div>
        <div class="columns three"><a href="?action=delete_all_from_search_engine&type={{.Type}}">Delete all from search engine</a></div>
        <div class="row">
          <div class="columns twelwe">
            <p>
              Total in DB: <b>{{ .TotalInDB }}</b>,
              Showing: <b>{{ .CurrentPageLabel }}</b>,
              Total in search engine: <b>{{ .TotalInSearchEngine }}</b>    
            </p>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="twelwe columns">
          {{$jobs := "jobs"}}
          {{$freelancers := "freelancers"}}
          <a class="{{if eq .Type $jobs}}button button-primary{{end}}" href="?type=jobs">Jobs</a>
          <a class="{{if eq .Type $freelancers}}button button-primary{{end}}" href="?type=freelancers">Freelancers</a>
          {{if ne (len .Entities) 0 }}
            <p>
              <a href="{{ .PrevPageURL }}&type={{.Type}}">{{ .PrevPageLabel }}</a>
              <a href="?type={{.Type}}">reset</a>
              <a href="{{ .NextPageURL }}&type={{.Type}}">{{ .NextPageLabel }}</a>
            </p>
          {{end}}
          <table class="u-full-width">
            <thead>
              <tr>
                <th>id</th>
                <th>name</th>
                <th>updatedAt</th>
                <th>action</th>
              </tr>
            </thead>
            <tbody>
              {{$page := .}}
              {{ range $docID, $value := .Entities }}
                <tr>
                  <td>{{ $docID }}</td>
                  <td>{{ $page.GetName $value }}</td>
                  <td>{{ $page.FormatTime $value.UpdatedAt }}</td>
                  <td>
                    {{if ne (len $page.Document) 0 }}
                      {{if eq $page.Document.id $docID }}
                        <a href="{{ $page.CurrentPageURL }}&type={{$page.Type}}">hide</a>
                      {{else}}
                        <a href="{{ $page.CurrentPageURL }}&action=get&docID={{$docID}}&type={{$page.Type}}">get</a>
                      {{end}}
                    {{else}}
                      <a href="{{ $page.CurrentPageURL }}&action=get&docID={{$docID}}&type={{$page.Type}}">get</a>
                    {{end}}
                    <a href="{{ $page.CurrentPageURL }}&action=import&docID={{$docID}}&type={{$page.Type}}">import</a>
                    <a href="{{ $page.CurrentPageURL }}&action=remove&docID={{$docID}}&type={{$page.Type}}">remove</a>
                  </td>
                </tr>
                {{if ne (len $page.Document) 0 }}
                  <tr>
                    {{if eq $page.Document.id $docID }}
                      <td colspan="4">
                        <table class="u-full-width">
                          <thead></thead>
                          <tbody>
                            {{ range $fieldName, $value := $page.Document.fields }}
                              <tr><td><b>{{$fieldName}}</b></td><td>{{ $value }}</td></tr>
                            {{ end }}
                          </tbody>
                        </table>
                      </td>
                    {{end}}
                  </tr>
                {{end}}
              {{ end }}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </body>
</html>`
