// The addModal object is a singleton web component for soliciting input in
// order to create a new row of data
//
// NB: addModal is very similar to editModal.  We have them as different
//     objects for conceptual clarity, but you could overload one modal to do
//     both.
var addModal = {
    
    // track if we've loaded and configured this modal yet
    init : false,

    // reset and show the modal
    Show : function() {
        // inject the modal into the DOM if we've never displayed it before
        if (!addModal.init) {
            $("body").append(hb_t.addModal());
            // wire up the buttons, and set the behavior on modal show
            $("#addModal").on("shown.bs.modal", function() {
                $("#addModal-smallnote").focus()
            });
            $("#addModal-save").click(addModal.save);
            $("#addModal-cancel").click(addModal.cancel);
            addModal.init = true;
        }
        // clear modal content and show it
        $("#addModal input").val("");
        $("#addModal").modal("show");
    },

    // This runs when the "cancel" button of the modal is clicked
    cancel : function() {
        $("#addModal").modal("hide");
    },
    
    // Send the new content to the server, and on success, update the page
    save : function() {
        // gather data from "required" fields
        var newdata = {
            smallnote : $('#addModal-smallnote').val()+"",
            bignote : $('#addModal-bignote').val()+"",
            favint : parseInt($('#addModal-favint').val()),
            favfloat : parseFloat($('#addModal-favfloat').val()),
        }
        // stop now if any data missing
        if (newdata.smallnote+"" === "" || newdata.bignote === "" ||
            newdata.favint+"" === "" || newdata.favfloat+"" === "")
        {
            window.alert("Missing field");
            return;
        }
        // get the optional field:
        if ($('#addModal-trickfloat').val() + "" !== "")
            newdata.trickfloat = parseFloat( $('#addModal-trickfloat').val());
        // send everything to the server, re-load the page on success
        $.ajax({
            type: "POST",
            url: "/data",
            contentType: "application/json",
            data: JSON.stringify(newdata),
            success : function(data) {
                dataList.DisplayLatest();
                $("#addModal").modal("hide");                
            },
            error : function(data) {
                window.alert("Unspecified error: " + data)
                $("#addModal").modal("hide");
            }
        });
    }
}
