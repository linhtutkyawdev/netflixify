const urlParams = new URLSearchParams(window.location.search);
const s = urlParams.get('s')?.toLowerCase();
const search = document.getElementById('s');

const main = document.getElementById('main');

let posts = JSON.parse(main.getAttribute('data-posts'));

if (s) {
  search.defaultValue = s;
  posts = posts.filter(
    (p) =>
      p.Tags.toLowerCase().split(' ').join('').includes(s) ||
      p.Title.toLowerCase().includes(s) ||
      p.Description.toLowerCase().includes(s)
  );
}

if (posts.length < 1) {
  main.innerHTML = 'No post found!';
}

const botUrl = main.getAttribute('data-botUrl');

document.getElementById(
  'banner__link'
).href = `${botUrl}?start=${posts[0].Video_id.slice(15)}`;
document.getElementById('banner__img').src = posts[0].G_thumbnail_path;

document.getElementById('new__container').innerHTML = posts
  .map(
    (p, index) => `
    <article class="card__article swiper-slide">
      <button onclick="modal_${index}.show()" class="card__link">
        <img src="${p.Thumbnail_path}" alt="image" class="card__img"/>
        <div class="card__shadow"></div>
        <div class="card__data">
            <h3 class="card__name">${p.Title}</h3>
            <span class="card__category">Rating : ${p.Rating}%</span>
        </div>
        <i class="ri-heart-3-line card__like"></i>
      </button>
    </article>
      `
  )
  .join('');

document.getElementById('modal__container').innerHTML = posts
  .map(
    (p, index) => `				
    <dialog id="modal_${index}" class="fixed top-0 left-0 w-full h-full bg-slate-900/20 backdrop-blur-lg z-10">
    <div class="container md:pl-80 md:pt-32 pt-40 px-8 text-white">
          <div class="flex md:flex-row flex-col justify-start items-center mb-4">
            <img alt="image" src="${
              p.G_thumbnail_path
            }" class="md:w-1/2 w-full rounded-lg"/>
            <div class="flex flex-row-reverse md:flex-col">
            <a href="${botUrl}?start=${p.Video_id.slice(15)}">
              <button class="bg-blue-400 w-24 p-2 m-4 text-white flex items-center rounded-md">Watch 🍿</button>
            </a>
            <form method="dialog">
              <button class="bg-rose-500 w-24 p-2 m-4 text-white flex items-center rounded-md">Close <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-circle-x ml-2 w-4 h-4"><circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg></button>
            </form>
            </div>
          </div>
          <p>${p.Description}</p>
    </div>
  </dialog>
      `
  )
  .join('');

const moviesIndexes = posts
  .map((p, index) =>
    p.Tags.toLowerCase().split(' ').join('').includes('movie') ? index : -1
  )
  .filter((p) => p !== -1);

const seriesIndexes = posts
  .map((p, index) =>
    p.Tags.toLowerCase().split(' ').join('').includes('serie') ? index : -1
  )
  .filter((p) => p !== -1);

if (moviesIndexes.length < 1 || s === 'movie') {
  document.getElementById('movies__list').hidden = true;
} else {
  document.getElementById('movies__container').innerHTML = moviesIndexes
    .map(
      (index) => `
    <article class="card__article swiper-slide">
      <button onclick="modal_${index}.show()" class="card__link">
        <img src="${posts[index].Thumbnail_path}" alt="image" class="card__img"/>
        <div class="card__shadow"></div>
        <div class="card__data">
            <h3 class="card__name">${posts[index].Title}</h3>
            <span class="card__category">Rating : ${posts[index].Rating}%</span>
        </div>
        <i class="ri-heart-3-line card__like"></i>
      </button>
    </article>
      `
    )
    .join('');
}

if (seriesIndexes.length < 1 || s === 'serie') {
  document.getElementById('series__list').hidden = true;
} else {
  document.getElementById('series__container').innerHTML = seriesIndexes
    .map(
      (index) => `
  <article class="card__article swiper-slide">
    <button onclick="modal_${index}.show()" class="card__link">
      <img src="${posts[index].Thumbnail_path}" alt="image" class="card__img"/>
      <div class="card__shadow"></div>
      <div class="card__data">
          <h3 class="card__name">${posts[index].Title}</h3>
          <span class="card__category">Rating : ${posts[index].Rating}%</span>
      </div>
      <i class="ri-heart-3-line card__like"></i>
    </button>
  </article>
    `
    )
    .join('');
}
