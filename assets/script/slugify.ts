/**
 * @description Slugify the input value and show it in another input box.
 */
class Slugify {
  constructor(slugEl?: HTMLInputElement) {
      if (!slugEl) {
          return;
      } else if (!(slugEl instanceof HTMLElement)) {
          slugEl = document.querySelector(slugEl);
      } else if (!(slugEl instanceof HTMLInputElement)) {
          return;
      }

      const targetSelector = slugEl.getAttribute("data-target");
      if (!targetSelector) {
          return;
      }

      const targetEl = document.querySelector<HTMLInputElement>(targetSelector);
      if (!targetEl || !(targetEl instanceof HTMLInputElement)) {
          return;
      }

      slugEl.addEventListener("input", () => {
          targetEl.value = slugEl.value.replace(/[-\s]+/g, '-').toLowerCase();
      });
  }

  static init(el?: HTMLElement) {
      if (!el) {
          el = document.body;
      } else if (!(el instanceof HTMLElement)) {
          el = document.querySelector(el);
      }

      const slugEls = el.querySelectorAll<HTMLInputElement>(`[data-component="slugify"]`);
      const slugs = [];

      for (let i = 0; i < slugEls.length; i++) {
          slugs.push(new Slugify(slugEls[i]));
      }

      return slugs;
  }
}

export default Slugify;
