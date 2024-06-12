
const bgImg = document.getElementById("bgImg");
const form = document.getElementById("thumbnail-form");
const [imgFile, _, imgSrc, title, subtitle, categories] = form;
bgImg.addEventListener("click", (e) => e.target.id == "bgImg" && imgFile.click());

imgFile.addEventListener("change", (e) => {
    const [file] = imgFile.files
    if (file) {
        bgImg.src = URL.createObjectURL(file);
    }
});

document.getElementById("saveBtn").addEventListener("click", (e) => {
    title.value = document.getElementById("title").innerText;
    subtitle.value = document.getElementById("subtitle").innerText;
    categories.value = document.getElementById("categories").innerText.split("\n").join(",");
    if (imgFile.files.length == 0) {
        imgSrc.value = bgImg.getAttribute("src-data");
    }
    e.target.innerHTML = "Loading.."
})

