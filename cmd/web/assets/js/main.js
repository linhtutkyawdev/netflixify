const urlParams = new URLSearchParams(window.location.search);
const s = urlParams.get('s')?.toLowerCase();
const search = document.getElementById('s');

// const t = urlParams.get('t').toLowerCase();

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

const movies = posts.filter((p) =>
  p.Tags.toLowerCase().split(' ').join('').includes('movie')
);

const series = posts.filter((p) =>
  p.Tags.toLowerCase().split(' ').join('').includes('serie')
);

if (movies.length < 1) {
  document.getElementById('movies__list').hidden = true;
}
if (series.length < 1) {
  document.getElementById('series__list').hidden = true;
}

document.getElementById(
  'banner__link'
).href = `${botUrl}?start=${posts[0].Video_id.slice(15)}`;
document.getElementById('banner__img').src = posts[0].G_thumbnail_path;

document.getElementById('new__container').innerHTML = posts
  .map(
    (p) => `
      <article class="card__article swiper-slide">
          <a href="${botUrl}?start=${p.Video_id.slice(15)}" class="card__link">
              <img src="${p.Thumbnail_path}" alt="image" class="card__img"/>
              <div class="card__shadow"></div>
              <div class="card__data">
                  <h3 class="card__name">${p.Title}</h3>
                  <span class="card__category">Rating : ${p.Rating}%</span>
              </div>
              <i class="ri-heart-3-line card__like"></i>
          </a>
      </article>
      `
  )
  .join('');
document.getElementById('movies__container').innerHTML = movies
  .map(
    (m) => `
    <article class="card__article swiper-slide">
        <a href="${botUrl}?start=${m.Video_id.slice(15)}" class="card__link">
			<img src="${m.Thumbnail_path}" alt="image" class="card__img"/>
			<div class="card__shadow"></div>
			<div class="card__data">
				<h3 class="card__name">${m.Title}</h3>
				<span class="card__category">Rating : ${m.Rating}%</span>
			</div>
			<i class="ri-heart-3-line card__like"></i>
		</a>
	</article>
    `
  )
  .join('');
document.getElementById('series__container').innerHTML = series
  .map(
    (s) => `
    <article class="card__article swiper-slide">
        <a href="${botUrl}?start=${s.Video_id.slice(15)}" class="card__link">
			<img src="${s.Thumbnail_path}" alt="image" class="card__img"/>
			<div class="card__shadow"></div>
			<div class="card__data">
				<h3 class="card__name">${s.Title}</h3>
				<span class="card__category">Rating : ${s.Rating}%</span>
			</div>
			<i class="ri-heart-3-line card__like"></i>
		</a>
	</article>
    `
  )
  .join('');
