/**
 * Created by kanta on 1/21/17.
 */

document.addEventListener("DOMContentLoaded", function () {
    // data = [
    //    []
    // ];
    makeAjaxCall(null,"redistest",null,initHandsonTable,null);

});


function initHandsonTable(data, textStatus, jqXHR){
    /*
     var data = [
     ["", "Ford", "Volvo", "Toyota", "Honda"],
     ["2016", 10, 11, 12, 13],
     ["2017", 20, 11, 14, 13],
     ["2018", 30, 15, 12, 13]
     ],
     */

    var container = document.getElementById("example"),
        searchField = document.getElementById('search_field'),
        resultCount = document.getElementById('resultCount'),
        hot,
        searchResultCount = 0;

    //Search
    function searchResultCounter(instance, row, col, value, result) {
        Handsontable.Search.DEFAULT_CALLBACK.apply(this, arguments);

        if (result) {
            searchResultCount++;
        }
    }

    hot = new Handsontable(container, {
        afterChange: function (changes, source) {
            if(!changes || changes.length == 0){
                console.log("Changes is null.");
                return;
            }
            console.log("This is fired after change.");
            console.log("CHANGES: ", changes);
            console.log("SOURCE: ", source);
            var data = formatChangesForRedis(changes);
            makeAjaxCall("POST", "setCells", data, null,null);
        },
        //afterChange se triggera kada se dogodi promjena.
        afterCreateRow: function (index, amount, source) {
            console.log("index: ", index);
            console.log("amount: ", amount);
            console.log("SOURCE: ", source);
        },
        search: {
            callback: searchResultCounter
        },
        data: data,
        //minSpareCols: 1,
        //minSpareRows: 1,
        rowHeaders: true,
        colHeaders: true,
        contextMenu: true,
        manualRowResize: true,
        columnSorting: true,
        manualColumnResize: true,
        maxRows: 7,
        maxCols: 7
    });

    Handsontable.Dom.addEvent(searchField, 'keyup', function(event) {
        var queryResult;

        searchResultCount = 0;
        queryResult = hot.search.query(this.value);
        console.log("queryResult", queryResult);
        resultCount.innerText = searchResultCount.toString();
        hot.render();
    });

    // This will be usefull for getting data dump from the whole table.
    /*
     function bindDumpButton() {
     if (typeof Handsontable === "undefined") {
     return;
     }

     Handsontable.Dom.addEvent(document.body, 'click', function (e) {

     var element = e.target || e.srcElement;

     if (element.nodeName == "BUTTON" && element.name == 'dump') {
     var name = element.getAttribute('data-dump');
     var instance = element.getAttribute('data-instance');
     var hot = window[instance];
     console.log('data of ' + name, hot.getData());
     }
     });
     }
     bindDumpButton();
     */
}

function formatChangesForRedis(changes){
    if(!changes || changes.length == 0){
        console.log("Changes should be array of at least one element.")
        return;
    }
    var data = [];

    if(!(changes[0] instanceof Array) && isNumber(changes[0])){
        data.push(_getCellData(changes));
    }

    for(var i = 0; i < changes.length; i++){
        data.push(_getCellData(changes[i]));
    }

    return data;
}

function _getCellData(cell){
    return {
        "row": cell[0], //row
        "column": cell[1], //column
        "value": cell[3], //new value
        "redisKey": "CELL_"+cell[0].toString()+"_"+cell[1].toString()
    }
}

function isNumber(obj) { return !isNaN(parseFloat(obj)) }

function makeAjaxCall(type, url, data, successCallback, dataType){
    if(!url){
        console.log("Can't make ajax call without providing url.");
        return;
    }
    type = type || "GET";
    dataType = dataType || "json";
    data = data || {};
    successCallback = successCallback || function(data, textStatus, jqXHR){
            console.log("Response data: ", data);
            console.log("Text status: ", textStatus);
            console.log("JQXHR: ", jqXHR);
        };
    var ajaxOptions = {
        type: type,
        url: url,
        data: data,
        success: successCallback,
        dataType: dataType
    };

    if (type !== "GET"){
        ajaxOptions["contentType"] = "application/json; charset=utf-8";
        ajaxOptions["data"] = JSON.stringify(data);
    }

    $.ajax(ajaxOptions);
}