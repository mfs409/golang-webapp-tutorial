{{! This file provides a Bootstrap modal for allowing the user to edit a}}
{{! row of data}}
<div id="editModal" class="modal fade" role="dialog">
 <div class="modal-dialog">
  <div class="modal-content">
   <div class="modal-header">
    <button type="button" class="close" data-dismiss="modal">
     <span class="glyphicon glyphicon-remove"></span>
    </button>
    <h4 id="modal-header" class="modal-title">Edit Data</h4>
   </div>
   <div class="modal-body">
    <div class="form-group">
     <label for="editModal-smallnote">Small Note</label>
     <input type="text" class="form-control" id="editModal-smallnote" placeholder="-" />
    </div>
    <div class="form-group">
     <label for="editModal-bignote">Big Note</label>
     <input type="text" class="form-control" id="editModal-bignote" placeholder="-" />
    </div>
    <div class="form-group">
     <label for="editModal-favint">Favorite Integer</label>
     <input type="text" class="form-control" id="editModal-favint" placeholder="-" />
    </div>
    <div class="form-group">
     <label for="editModal-favfloat">Favorite Float</label>
     <input type="text" class="form-control" id="editModal-favfloat" placeholder="-" />
    </div>
    <div class="form-group">
     <label for="editModal-trickfloat">Tricky Float</label>
     <input type="text" class="form-control" id="editModal-trickfloat" placeholder="-" />
    </div>
   </div>
   <div class="modal-footer">
    <button id="editModal-save" type="button" class="btn btn-success">Save</button>
    <button id="editModal-cancel" type="button" class="btn btn-danger">Close</button>
   </div>
   <input type="hidden" id="editModal-id" />
  </div>
 </div>
</div>
