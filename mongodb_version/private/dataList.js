// The dataList object is a singleton web component for presenting all of the
// data currently in the database.
//
// NB: The component is very simple... on any change to the data, we're going
//     to re-fetch the latest data and regenerate the entire component.  It
//     would be straightforward to extend this component so that it could
//     update individual rows, but in the interest of keeping the tutorial
//     simple, we won't.

var dataList = {

    // Fetch the latest content, re-create the component, and (re)place it in
    // the DOM.
    DisplayLatest : function() {
        $.ajax({
            url: "/data",
            success: function(data) {
                // Put the data into the DOM
                $(".container").html(hb_t.dataList({d:data}));
                
                // wire up click handlers for buttons
                //
                // NB: there is one add button, but potentially many edit and
                //     delete buttons
                $(".dataList-addbtn").click(addModal.Show);
                $(".dataList-deletebtn").click(dataList.deleteRow);
                $(".dataList-editbtn").click(dataList.editRow);
            },
            dataType: "json"
        })        
    },

    // when a 'delete' button is clicked, we use the data-id of the
    // button to know what row to delete
    deleteRow : function() {
            var id = $(this).data("id");
            $.ajax({
                type: "DELETE",
                url: "/data/"+id,
                contentType: "application/json",
                success : function(data) {
                    dataList.DisplayLatest();
                    location.reload();
                },
                error : function(data) {
                    window.alert("Unspecified error: " + data)
                }
            });

    },

    // when an 'edit' button is clicked, we need to extract the data from the
    // row, and then use it to populate the edit modal
    editRow: function() {
        var id = $(this).data("id");
        // find the row that holds this button, get its cells
        var tds = $(this).closest("tr").find("td");
        // pull text from the cells to produce the default values for the
        // modal, then show the modal with the data.
        var content = {
            id : id,
            smallnote : $(tds[0]).text(),
            bignote : $(tds[1]).text(),
            favint : $(tds[2]).text(),
            favfloat : $(tds[3]).text(),
            trickfloat : $(tds[4]).text()
        };
        editModal.Show(content);
    }
}
