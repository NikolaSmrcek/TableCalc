/**
 * Created by kanta on 1/21/17.
 */

var hot = null;

var stats = null,
    punctuations = null,
    wordsOfInterest = null,
    wordsOfNoInterest = null;

var handSonTables = {};

document.addEventListener("DOMContentLoaded", function () {
    // data = [
    //    []
    // ];
    makeAjaxCall(null,"getInitData",null,initHandsonTable,null);
    makeAjaxCall(null,"getMapReduceResult",null,getMapReduceResult,null);
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
            var data = formatChangesForRedis(changes);
            makeAjaxCall("POST", "setCells", data, null,null);
        },
        //afterChange se triggera kada se dogodi promjena.
        afterCreateRow: function (index, amount) {
            console.log("index: ", index);
            console.log("amount: ", amount);
            var data = formatHotGetDataForRedis(hot.getData());
            makeAjaxCall("POST", "setCells", data, null,null);
            makeAjaxCall("POST", "addToList", {"name": "rows", "value": index.toString()}, null,null);
        },
        afterCreateCol: function(index, amount) {
            console.log("index: ", index);
            console.log("amount: ", amount);
            var data = formatHotGetDataForRedis(hot.getData());
            makeAjaxCall("POST", "setCells", data, null,null);
            makeAjaxCall("POST", "addToList", {"name": "columns", "value": index.toString()}, null,null);
        },
        afterRemoveRow: function(index, amount) {
            var data = formatHotGetDataForRedis(hot.getData());
            makeAjaxCall("POST", "setCells", data, null,null);
            makeAjaxCall("POST", "removeFromList", {"name": "rows"}, null,null);
        },
        afterRemoveCol: function(index, amount) {
            var data = formatHotGetDataForRedis(hot.getData());
            makeAjaxCall("POST", "setCells", data, null,null);
            makeAjaxCall("POST", "removeFromList", {"name": "columns"}, null,null);
        },
        afterColumnSort: function(column, order) {
            var data = formatHotGetDataForRedis(hot.getData());
            makeAjaxCall("POST", "setCells", data, null,null);
        },
        beforeChange: function(changes, source) {
            //TODO implement SUM and AVG
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

function formatHotGetDataForRedis(all_data){
    if(!all_data || all_data.length == 0){
        console.log("Changes should be array of at least one element.");
        return;
    }
    var data = [];

    for(var i = 0; i < all_data.length; i++){
        //getting array of the rows columns so the cell value will be row[j]
        var row = all_data[i];
        for(var j = 0; j < row.length; j++){
            var cell = {
                "row": i,
                "column": j,
                "value": row[j]
            };
            cell["redisKey"] = "CELL_"+i.toString()+"_"+j.toString();
            data.push(cell);
        }
    }

    return data;
}

function formatChangesForRedis(changes){
    if(!changes || changes.length == 0){
        console.log("Changes should be array of at least one element.");
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

//adding listener for submit
$("form#data").submit(function(e){
    if(document.getElementById("inputFile").files.length == 0){
        alert("Select at least one file to upload.");
        return false;
    }
    var formData = new FormData($(this)[0]);

    $.ajax({
        url: "upload",
        type: 'POST',
        data: formData,
        async: false,
        success: getMapReduceResult,
        cache: false,
        contentType: false,
        processData: false
    });

    return false;
});

function getMapReduceResult(data){
    if (data) {
        _loadHanson("stats", extractFromObjectForHandsonTable(data.stats));
        _loadHighChartsGraph(extractDataForHighChartsGraph(data.stats), "statsGraph", "Overall words statistics", "Overall statistics");

        _loadHanson("punctuations", extractFromObjectForHandsonTable(data.punctuations));
        _loadHighChartsGraph(extractDataForHighChartsGraph(data.punctuations), "punctuationsGraph", "Top 5 punctuations", "Punctuations");

        _loadHanson("wordsOfInterest", extractFromObjectForHandsonTable(data.wordsOfInterest));
        _loadHighChartsGraph(extractDataForHighChartsGraph(data.wordsOfInterest), "wordsOfInterestGraph", "Top 5 words of interest", "wordsOfInterest");

        _loadHanson("wordsOfNoInterest", extractFromObjectForHandsonTable(data.wordsOfNoInterest));
        _loadHighChartsGraph(extractDataForHighChartsGraph(data.wordsOfNoInterest), "wordsOfNoInterestGraph", "Top 5 words of no interest", "wordsOfNoInterest");
    }
}

function destoryHotTable(containerId){

}

function extractFromObjectForHandsonTable(data){
    handsonData = [];

    for (var key in data){
        if (data.hasOwnProperty(key)) {
            handsonData.push([key, data[key]]);
        }
    }
    console.log("Handson data: ", handsonData);
    return handsonData;
}

/*
 var data = [
 ["", "Ford", "Volvo", "Toyota", "Honda"],
 ["2016", 10, 11, 12, 13],
 ["2017", 20, 11, 14, 13],
 ["2018", 30, 15, 12, 13]
 ],
 */

function _loadHanson(containerId, data){
    var container = document.getElementById(containerId);
    console.log("Loading handson table.");
    if(handSonTables[containerId]){
        handSonTables[containerId].loadData(data);
    }
    else{
        handSonTables[containerId] = new Handsontable(container, {
            data: data,
            //minSpareCols: 1,
            //minSpareRows: 1,
            rowHeaders: true,
            colHeaders: true,
            contextMenu: true,
            manualRowResize: true,
            columnSorting: true,
            manualColumnResize: true,
            readOnly: true
        });
    }
}

function extractDataForHighChartsGraph(data){
    highChartsdata = [];

    for (var key in data){
        if (data.hasOwnProperty(key)) {
            highChartsdata.push({"name": key, "y":parseInt(data[key], 10)} );
            //highChartsdata.push([key, data[key]]);
        }
    }

    highChartsdata.sort(function(a, b){
       return  b.y - a.y;
    });

    highChartsdata.splice(5);

    return highChartsdata;
}

function _loadHighChartsGraph(data, container, text, seriesName){
    /*
     [{
     name: 'Microsoft Internet Explorer',
     y: 56.33
     }, {
     name: 'Chrome',
     y: 24.03,
     sliced: true,
     selected: true
     }, {
     name: 'Firefox',
     y: 10.38
     }, {
     name: 'Safari',
     y: 4.77
     }, {
     name: 'Opera',
     y: 0.91
     }, {
     name: 'Proprietary or Undetectable',
     y: 0.2
     }]
     */
    console.log("Highcharts data: ", data);
    Highcharts.chart(container, {
        chart: {
            plotBackgroundColor: null,
            plotBorderWidth: null,
            plotShadow: false,
            type: 'pie'
        },
        title: {
            text: text
        },
        tooltip: {
            //{point.percentage:.1f}
            pointFormat: '{series.name}: <b>{point.y}</b>'
        },
        plotOptions: {
            pie: {
                allowPointSelect: true,
                cursor: 'pointer',
                dataLabels: {
                    enabled: false
                },
                showInLegend: true
            }
        },
        series: [{
            name: seriesName,
            colorByPoint: true,
            data: data
        }]
    });
}