/**
 * @description Click a button to load an image from url entered into an input box.
 * Used by views/promo/new-banner.html
 */
class PreviewImage {
  constructor(previewEl: HTMLElement) {
    const sourceSelector = previewEl.getAttribute("data-image-source");
    if (!sourceSelector) {
      return;
    }

    const sourceEl = document.querySelector<HTMLInputElement>(sourceSelector);
    if (!sourceEl || !(sourceEl instanceof HTMLInputElement)) {
      return;
    }

    const targetSelector = previewEl.getAttribute("data-target");
    if (!targetSelector) {
      return;
    }
    const targetEl = document.querySelector<HTMLElement>(targetSelector);
    if (!targetEl || !(targetEl instanceof HTMLImageElement)) {
      return;
    }

    previewEl.addEventListener("click", e => {
      const url = sourceEl.value;
      targetEl.src = url;
    });
  }

  static init(el?: HTMLElement) {
    if (!el) {
      el = document.body;
    } else if (!(el instanceof HTMLElement)) {
      el = document.querySelector(el);
    }

    const previewEls = el.querySelectorAll<HTMLElement>('[data-component="preview-image"]');
    const previews = [];

    for (let i = 0; i < previewEls.length; i++) {
      previews.push(new PreviewImage(previewEls[i]));
    }

    return previews;
  }
}

export default PreviewImage;
