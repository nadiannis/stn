const linkList = document.querySelector('.link-list');
const linkDeleteButtons = document.querySelectorAll('.btn-link-delete');

linkDeleteButtons.forEach((linkDeleteButton) => {
  const id = linkDeleteButton.dataset.id;
  linkDeleteButton.addEventListener('click', () => deleteLink(id));
});

linkList.addEventListener('click', (e) => {
  if (e.target.classList.contains('btn-short-link-copy')) {
    copyShortLink(e);
  }
});
