package importer

var htmlTemplate = `
<html>
  <head>
    <title>Fairlance Importer</title>
    <link href="http://getskeleton.com/dist/css/normalize.css" rel="stylesheet" type="text/css"/>
    <link href="http://getskeleton.com/dist/css/skeleton.css" rel="stylesheet" type="text/css"/>
    <style>
      html {
          overflow-y: scroll; 
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="row">
        <div class="column one"><a href="?type={{.Type}}">Reload</a></div>
        <div class="columns eleven">
          <b>Status: {{.Message}}</b>
        </div>
      </div>
      <div class="row">
        <div class="twelwe columns">
          {{$jobs := "jobs"}}
          {{$freelancers := "freelancers"}}
          <a class="button {{if eq .Type $jobs}}button-primary{{end}}" href="?type=jobs">Jobs</a>
          <a class="button {{if eq .Type $freelancers}}button-primary{{end}}" href="?type=freelancers">Freelancers</a>
        </div>
      </div>
      <div class="row">
        <div class="twelwe columns">
          <a class="button" href="?action=re_generate_test_data&type={{.Type}}">Re-generate test data</a>
          <a class="button" href="?action=import_all&type={{.Type}}">Import to search engine</a>
          <a class="button" href="?action=delete_all_from_db&type={{.Type}}" onclick="return confirm('Are you sure?')">Delete from DB</a>
          <a class="button" href="?action=delete_all_from_search_engine&type={{.Type}}" onclick="return confirm('Are you sure?')">Delete from search engine</a>
        </div>
      </div>
      <div class="row">
        <div class="twelwe columns">
          <div class="row">
            <div class="columns twelwe">
              <p>
                {{if ne (len .Entities) 0 }}
                  {{if or (ne .PrevPageURL "") (ne .NextPageURL "")  }}
                    <a href="?type={{.Type}}">reset</a> |
                  {{ end }}
                  <a href="{{ .PrevPageURL }}&type={{.Type}}">{{ .PrevPageLabel }}</a>
                  <b>{{ .CurrentPageLabel }}</b>
                  <a href="{{ .NextPageURL }}&type={{.Type}}">{{ .NextPageLabel }}</a> |
                {{end}}
                Total in DB: <b>{{ .TotalInDB }}</b>,
                Total in search engine: <b>{{ .TotalInSearchEngine }}</b>
              </p>
            </div>
          </div>
        </div>
        <div class="row">
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
                  <td title="{{ $page.FormatTime $value.UpdatedAt }}">{{ $page.FormatTimeHuman $value.UpdatedAt }}</td>
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
