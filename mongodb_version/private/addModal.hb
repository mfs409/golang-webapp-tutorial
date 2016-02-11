{{! This file provides a Bootstrap modal for allowing the user to add a new}}
{{! row of data}}
<div id="addModal" class="modal fade" role="dialog">
 <div class="modal-dialog">
  <div class="modal-content">
   <div class="modal-header">
    <button type="button" class="close" data-dismiss="modal">
     <span class="glyphicon glyphicon-remove"></span>
    </button>
    <h4 id="modal-header" class="modal-title">Add New Data</h4>
   </div>
   <div class="modal-body">
    <div class="form-group">
     <label for="addModal-smallnote">Small Note</label>
     <input type="text" class="form-control" id="addModal-smallnote" placeholder="-">
    </div>
    <div class="form-group">
     <label for="addModal-bignote">Big Note</label>
     <input type="text" class="form-control" id="addModal-bignote" placeholder="-">
    </div>
    <div class="form-group">
     <label for="addModal-favint">Favorite Integer</label>
     <input type="text" class="form-control" id="addModal-favint" placeholder="-">
    </div>
    <div class="form-group">
     <label for="addModal-favfloat">Favorite Float</label>
     <input type="text" class="form-control" id="addModal-favfloat" placeholder="-">
    </div>
    <div class="form-group">
     <label for="addModal-trickfloat">Tricky Float</label>
     <input type="text" class="form-control" id="addModal-trickfloat" placeholder="-">
    </div>
   </div>
   <div class="modal-footer">
    <button id="addModal-save" type="button" class="btn btn-success">Save</button>
    <button id="addModal-cancel" type="button" class="btn btn-danger">Close</button>
   </div>
  </div>
 </div>
</div>
