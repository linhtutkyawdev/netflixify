
const template = document.getElementById("thumbnail-template");
const imgSrc = document.getElementsByName("imgSrc")[0];
const title = document.getElementById("title");
const downloadBtn = document.getElementById("downloadBtn");
const thumbnailContainer = document.getElementById("thumbnail-container");

template.addEventListener("click", (e) => e.target.id == "thumbnail-template" && imgSrc.click());

imgSrc.addEventListener("change", (e) => {
    const [file] = imgSrc.files
    if (file) {
        template.style.backgroundImage = "url(" + URL.createObjectURL(file) + ")";
    }
});

downloadBtn.addEventListener("click", (e) => {
    e.preventDefault();
    alert(HOST)
    fetch(":8080/convert/html2image?u=kode&p=31545&url=http://localhost:3000/thumbnail&customClip=true&clipX=0&clipY=0&clipWidth=720&clipHeight=405&clipScale=1&format=png")
        .then((response) => {
            return response.blob()
        })
        .then((blob) => {
            console.log('Blob:', blob)
            const url = URL

                .createObjectURL(blob)

            const link = document
                .createElement('a')
            link
                .href = url
            link
                .download = title.textContent.toLowerCase().replace(/\s/g, '') + '.png'
            // The name for the downloaded file
            document
                .body
                .appendChild(link)
            link
                .click()
            document
                .body
                .removeChild(link)

            URL.revokeObjectURL(url)
        })
        .catch(console.error)
})

// let isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);

// if (isMobile) {
//   // User is accessing the page on a mobile device
//   console.log("Mobile device detected");
//   thumbnailContainer.classList.replace("w-[720px]", "w-screen");
//   form.classList.add("mx-auto");
// } else {
//   // User is accessing the page on a desktop device
//   console.log("Desktop device detected");
// }
