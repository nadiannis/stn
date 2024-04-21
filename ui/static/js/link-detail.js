const linkDeleteButton = document.querySelector('.btn-link-delete');
const shortLinkCopyButton = document.querySelector('.btn-short-link-copy');
const id = location.pathname.split('/')[2];

linkDeleteButton.addEventListener('click', () => deleteLink(id));
shortLinkCopyButton.addEventListener('click', copyShortLink);
