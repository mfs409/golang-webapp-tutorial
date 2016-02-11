{{! This file provides a template for generating a table of data.  We use a}}
{{! little bit of Bootstrap styling on the table.}}
<table id="dataList" class="table table-striped" id="dataTable">
 <thead>
  <tr>
   <th>Small Note</th>
   <th>Big Note</th>
   <th>Favorite Int</th>
   <th>Favorite Float</th>
   <th>Extra Float?</th>
   <th><a href='#editModal' data-toggle='modal' title='Add'
          class='btn btn-default btn-sm dataList-addbtn' data-id=''>
     <span class='glyphicon glyphicon-plus'></span></a>
   </th>
  </tr>
 </thead>
 <tbody>
  {{#each d}}
  <tr>
   <td>{{smallnote}}</td>
   <td>{{bignote}}</td>
   <td>{{favint}}</td>
   <td>{{favfloat}}</td>
   <td>{{trickfloat}}</td>
   <td>
    <a href="#" title="Edit" class="btn btn-default btn-sm dataList-editbtn"
       data-id="{{id}}">
     <span class="glyphicon glyphicon-pencil"></span>
    </a>
    <a href="#" title="Delete" class="btn btn-default btn-sm dataList-deletebtn"
       data-id="{{id}}">
     <span class="glyphicon glyphicon-remove"></span></a>
   </td>
  </tr>
  {{/each}}
 </tbody>
</table>

