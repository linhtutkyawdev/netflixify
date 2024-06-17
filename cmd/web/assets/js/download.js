const saveBtn = document.getElementById("saveBtn")
if (saveBtn != null) {
    saveBtn.addEventListener("click", (e) => {
        e.preventDefault()

        fetch(saveBtn.getAttribute("imgSrc-data"))
            .then((response) => {
                return response.blob()
            })
            .then((blob) => {
                const url = URL.createObjectURL(blob)
                const link = document.createElement('a')
                link.href = url
                link.download = saveBtn.getAttribute("fileName-data")
                document.body.appendChild(link)
                link.click()
                document.body.removeChild(link)
                URL.revokeObjectURL(url)
                saveBtn.innerHTML = "Download Again"
            })
            .catch(console.error);

        document.getElementById("thumbnail-container").contentEditable = false;
        document.getElementById("title").classList.remove("cursor-text");
        document.getElementById("subtitle").classList.remove("cursor-text");
        document.getElementById("categories").childNodes.forEach((n) => n.classList.remove("cursor-text"));
        bgImg.classList.remove("cursor-pointer");
        imgFile.disabled = true;
    })

    saveBtn.click();
}