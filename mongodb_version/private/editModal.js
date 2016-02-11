// The editModal object is a singleton web component for soliciting input in
// order to change the content of a row of data
//
// NB: editModal is very similar to addModal.  We have them as different
//     objects for conceptual clarity, but you could overload one modal to do
//     both.
var editModal = {
    
    // track if we've loaded and configured this modal yet
    init : false,

    // populate and show the modal
    Show : function(data) {
        // inject the modal into the DOM if we've never displayed it before
        if (!editModal.init) {
            $("body").append(hb_t.editModal());
            // wire up the buttons, and set the behavior on modal show
            $("#editModal").on("shown.bs.modal", function() {
                $("#editModal-smallnote").focus()
            });
            $("#editModal-save").click(editModal.save);
            $("#editModal-cancel").click(editModal.cancel);
            editModal.init = true;
        }
        // set modal content and show it
        $("#editModal-smallnote").val(data.smallnote);
        $("#editModal-bignote").val(data.bignote);
        $("#editModal-favint").val(data.favint);
        $("#editModal-favfloat").val(data.favfloat);
        $("#editModal-trickfloat").val(data.trickfloat);
        $("#editModal-id").val(data.id);
        $("#editModal").modal("show");
    },

    // This runs when the "cancel" button of the modal is clicked
    cancel : function() {
        $("#editModal").modal("hide");
    },
    
    // Send the new content to the server, and on success, update the page
    save : function() {
        // gather data from "required" fields
        var newdata = {
            smallnote : $('#editModal-smallnote').val()+"",
            bignote : $('#editModal-bignote').val()+"",
            favint : parseInt($('#editModal-favint').val()),
            favfloat : parseFloat($('#editModal-favfloat').val()),
            id : $("#editModal-id").val() // NB: shouldn't ever be blank
        }
        // stop now if any data missing
        if (newdata.smallnote+"" === "" || newdata.bignote === "" ||
            newdata.favint+"" === "" || newdata.favfloat+"" === "")
        {
            window.alert("Missing field");
            return;
        }
        // get the optional field:
        if ($('#editModal-trickfloat').val() + "" !== "")
            newdata.trickfloat = parseFloat( $('#editModal-trickfloat').val());
        // send everything to the server, re-load the page on success
        $.ajax({
            type: "PUT",
            url: "/data/"+newdata.id,
            contentType: "application/json",
            data: JSON.stringify(newdata),
            success : function(data) {
                dataList.DisplayLatest();
                $("#editModal").modal("hide");                
            },
            error : function(data) {
                window.alert("Unspecified error: " + data)
                $("#editModal").modal("hide");
            }
        });
    }
}
