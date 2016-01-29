// store all data from the server, so that we can easily populate the form
// when 'edit' is pressed
var origData;

// Request some code to run when the page is done loading
$(document).ready(function() {
    fetch();
});

// fetch the data from the server and use it to build the table
function fetch() {
    $.ajax({
        url: "/data",
        success: function(data) {
            origData = data; // save it, to avoid server accesses later
            for (var i = 0; i < data.length; ++i) {
                makeRow(i);
            }
        },
        dataType: "json"
    })
}

// build a row in the table, using an index into origData
function makeRow(i) {
    // add a row
    $('#dataTable').append('<tr id="row'+i+'"></tr>');
    // add cells for data
    $('#row'+i).append('<td id="smallnote-'+origData[i].id+'">'+origData[i].smallnote+'</td>')
    $('#row'+i).append('<td id="bignote-'+origData[i].id+'">'+origData[i].bignote+'</td>')
    $('#row'+i).append('<td id="favint-'+origData[i].id+'">'+origData[i].favint+'</td>')
    $('#row'+i).append('<td id="favfloat-'+origData[i].id+'">'+origData[i].favfloat+'</td>')
    $('#row'+i).append('<td id="trickfloat-'+origData[i].id+'">'+origData[i].trickfloat+'</td>')
    // add buttons for edit and delete
    $('#row'+i).append('<td>' + '<a href="#editModal" data-toggle="modal" title="Edit" class="btn btn-default btn-sm edit-btn" data-id="'+i+'"><span class="glyphicon glyphicon-pencil"></span></a>' + '<a href="#" title="Delete" class="btn btn-default btn-sm delete-btn" data-id="'+i+'"><span class="glyphicon glyphicon-remove"></span></a>' + '</td>');
}

// wire up a click handler for 'edit' buttons
$(document).on("click", ".edit-btn", function() {
    // populate the fields of the modal form
    var i = $(this).data("id");
    $('#id').val(origData[i].id)
    $('#smallnote').val(origData[i].smallnote)
    $('#bignote').val(origData[i].bignote)
    $('#favint').val(origData[i].favint)
    $('#favfloat').val(origData[i].favfloat)
    $('#trickfloat').val(origData[i].trickfloat)

    // set modal title
    $("#modal-header").text("Edit Entry")    
    
    // wire up the modal "save" button
    $("#modal-save").unbind();
    $("#modal-save").click(function(){
        // wrap up the new data in an object that is easy to send
        var newdata = {
            id : $('#id').val()+"",
            smallnote : $('#smallnote').val()+"",
            bignote : $('#bignote').val()+"",
            favint : parseInt($('#favint').val()),
            favfloat : parseFloat($('#favfloat').val()),
            trickfloat : parseFloat( $('#trickfloat').val()),
        }
        if ($('#trickfloat').val() + "" !== "")
            newdata.trickfloat = parseFloat( $('#trickfloat').val());
        $.ajax({
            type: "PUT",
            url: "/data/"+newdata.id,
            contentType: "application/json",
            data: JSON.stringify(newdata),
            success : function(data) {
                location.reload();
            },
            error : function(data) {
                window.alert("Unspecified error: " + data)
            }
        });
    });
});

// wire up a click handler for 'delete' buttons
$(document).on("click", ".delete-btn", function() {
    // get the id
    var i = $(this).data("id")
    var id = origData[i].id

    // send a DELETE
    $.ajax({
        type: "DELETE",
        url: "/data/"+id,
        contentType: "application/json",
        success : function(data) {
            location.reload();
        },
        error : function(data) {
            window.alert("Unspecified error: " + data)
        }
    });
});

// wire up a click handler for 'add' button
$(document).on("click", ".add-btn", function() {
    // populate the fields of the modal form
    var i = $(this).data("id");
    $('#id').val("")
    $('#smallnote').val("")
    $('#bignote').val("")
    $('#favint').val("")
    $('#favfloat').val("")
    $('#trickfloat').val("")

    // set the title
    $("#modal-header").text("Add Entry")
    
    // wire up the modal "save" button
    $("#modal-save").unbind();
    $("#modal-save").click(function(){
        // wrap up the new data in an object that is easy to send
        var newdata = {
            smallnote : $('#smallnote').val()+"",
            bignote : $('#bignote').val()+"",
            favint : parseInt($('#favint').val()),
            favfloat : parseFloat($('#favfloat').val()),
            trickfloat : parseFloat( $('#trickfloat').val()),
        }
        if ($('#trickfloat').val() + "" !== "")
            newdata.trickfloat = parseFloat( $('#trickfloat').val());
        $.ajax({
            type: "POST",
            url: "/data",
            contentType: "application/json",
            data: JSON.stringify(newdata),
            success : function(data) {
                location.reload();
            },
            error : function(data) {
                window.alert("Unspecified error: " + data)
            }
        });
    });
});
