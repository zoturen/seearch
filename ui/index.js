const search_input = document.getElementById("search-text");
const search_btn = document.getElementById("search-btn");

search_btn.addEventListener("click", async () => {
  const search_term = search_input.value;
  const search_res = await fetch(
    `http://localhost:8080/search?q=${search_term}`,
  );

  if (search_res.status === 200) {
    //console.log(await search_res.json());
    const json_res = await search_res.json();
    const search_result_div = document.getElementById("search-result");
    search_result_div.innerHTML = "";

    const heading = document.createElement("h2");
    heading.textContent = "Search Results";

    search_result_div.appendChild(heading);

    const result = json_res.result;

    if (result.length < 1) {
      const no_result_p = document.createElement("p");
      no_result_p.textContent = "Sorry no results were found...";
      search_result_div.appendChild(no_result_p);
      return;
    }

    const result_list = document.createElement("ul");
    search_result_div.appendChild(result_list);

    result.forEach((result) => {
      const list_item = document.createElement("li");
      list_item.textContent = result.Key;
      result_list.appendChild(list_item);
    });
  }
});
