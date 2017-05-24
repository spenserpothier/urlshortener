
$( document ).ready(function (){
    $(".tagsinput").select2({
	tags: [],
	tokenSeparators: [',', ' '],
	ajax: {
	    url: "https://localhost:8000/api/tags",
	    dataType: 'json',
	    delay: 250,
	    data: function (params) {
		return {
		    q: params.term
		};      
		
	    },

	    processResults: function (data, params) {
		return {
		    results: $.map(data, function (item) {
			return {
			    text: item.Title,
			    id: item.Id    
			}
		    })
		}
	    }
	}
    });
});
