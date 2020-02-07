/**
 * @description Get browser's default timezone.
 * Used by: views/promo/new-schedule.html
 */
class ZoneOffset {
  readonly inputEl: HTMLInputElement;

  constructor(inputEl?: HTMLInputElement) {
      if (!inputEl) {
          return;
      }
      
      if (typeof inputEl == "string") {
        this.inputEl = document.querySelector(inputEl);
      } else {
        this.inputEl = inputEl;
      }

      if (!(this.inputEl instanceof HTMLInputElement)) {
          return;
      }
  }

  showZone() {
    if (!this.inputEl) {
      throw new Error("Please specify which input box to show time zone offset");
    }
    this.inputEl.value = this.isoOffset();
  }

  /**
   * @description Turns a JavaScript Date's time zone offset in minutes to ISO designator or offset string.
   * @returns {string}
   * @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/getTimezoneOffset
   */
  isoOffset() {
      const offset = new Date().getTimezoneOffset();

      if (offset == 0) {
          return "Z"
      }

      // The sign of getTimezoneOffset is opposite to ISO timezone format.
      let sign;
      if (offset <= 0) {
          sign = "+";
      } else {
          sign = "-";
      }

      // Must use the absolute value of offset first.
      const hour = Math.floor(Math.abs(offset) / 60).toFixed();
      const minute = Math.abs(Math.abs(offset ) % 60).toFixed();

      return `${sign}${hour.padStart(2, "0")}:${minute.padStart(2, "0")}`
  }

  static init(el?: HTMLElement) {
      if (!el) {
          el = document.body;
      } else if (!(el instanceof HTMLElement)) {
          el = document.querySelector(el);
      }

      const inputEls = el.querySelectorAll<HTMLInputElement>(`[data-component="zone-offset"]`);
      const inst = [];

      for (let i = 0; i < inputEls.length; i++) {
        const zoneOffset = new ZoneOffset(inputEls[i]);
        zoneOffset.showZone();

        inst.push(zoneOffset);
      }

      return inst;
  }
}

export default ZoneOffset;
