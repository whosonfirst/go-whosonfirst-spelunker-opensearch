window.addEventListener("load", function load(event){

    
    var facets_wrapper = document.querySelector("#whosonfirst-facets");

    if (! facets_wrapper){
	console.log("NOPE");
	return;
    }

    var current_url = facets_wrapper.getAttribute("data-current-url");
    var facets_url = facets_wrapper.getAttribute("data-facets-url");    

    if ((! current_url) || (! facets_url)){
	return;
    }

    var draw_facets = function(rsp){

	var f = rsp.facet.property;
	var id = "#whosonfirst-facets-" + f;

	var el = document.querySelector(id);

	if (! el){
	    console.log("Unable to find facet wrapper", id);
	    return;
	}

	var f_label = f;

	if (f == "iscurrent") {
	    f_label = "is current";
	}
	
	var label = document.createElement("h3");
	label.appendChild(document.createTextNode(f_label));
	
	var ul = document.createElement("ul");
	ul.setAttribute("class", "whosonfirst-facets");
	
	var results = rsp.results;
	var count = results.length;

	for (var i=0; i < count; i++){

	    var k = results[i].key;

	    if (k == ""){

		var sp = document.createElement("span");
		sp.setAttribute("class", "hey-look");
		sp.appendChild(document.createTextNode("undefined"));

		var sm = document.createElement("small");
		sm.appendChild(document.createTextNode(results[i].count));
		
		var item = document.createElement("li");
		item.appendChild(sp);
		item.appendChild(sm);

	    } else {

		var k_label = k;

		if (f == "current"){

		    switch (parseInt(k)){
			case 0:
			    k_label = "not current";
			    break;
			case -1:
			    k_label = "unknown";
			    break;
			default:
			    k_label = "current";
			    break;
		    }
			    
		}
		
		// Something something something is location.href really safe?
		// https://developer.mozilla.org/en-US/docs/Web/API/URL/URL

		var u = new URL(current_url, location.href);
		u.searchParams.set(f, k)

		var a = document.createElement("a");
		
		a.setAttribute("href", u.toString());
		a.setAttribute("class", "hey-look");
		a.appendChild(document.createTextNode(k_label));
		
		var sm = document.createElement("small");
		sm.appendChild(document.createTextNode(results[i].count));
		
		var item = document.createElement("li");
		item.appendChild(a);
		item.appendChild(sm);
	    }
	    
	    ul.appendChild(item);
	}

	el.appendChild(label);
	el.appendChild(ul);
    };
    
    var fetch_facet = function(f){

	// var url = facets_url + "?&facet=" + f;

	// Something something something is location.href really safe?
	// https://developer.mozilla.org/en-US/docs/Web/API/URL/URL
	
	var u = new URL(facets_url, location.href)
	u.searchParams.set("facet", f);
	var url = u.toString();
	
	fetch(url)
	    .then((rsp) => rsp.json())
	    .then((data) => {

		var count = data.length;

		for (var i=0; i < count; i++){
		    draw_facets(data[i]);
		}
		
	    }).catch((err) => {
		console.log("SAD", f, err);
	    });
    };
    
    var facets = facets_wrapper.getAttribute("data-facets");
    facets = facets.split(",");

    var count_facets = facets.length;

    for (var i=0; i < count_facets; i++){

	var f = facets[i];
	
	var el = document.createElement("div");
	el.setAttribute("id", "whosonfirst-facets-" + f);
	facets_wrapper.appendChild(el);

	fetch_facet(f);
    }

});
