import { VsPeriod } from "./period.js";
import { VsTimeline } from "./timeline.js";

const template = document.createElement("template");
template.innerHTML = `
  <style>
    .main-container {
      width: 500px;
      height: 1200px;
      display: grid;
      grid-template-columns: 1fr 1fr 1fr;
    }

    .timeline-container {
      width: 10rem; 
      position: relative;
    }

    .graph-container {
      width: 500px;
      height: 100%;
      position: relative;
      background-color: lightgray;
    }
  </style>

  <div class="main-container">
    <div class="timeline-container">
      <slot name="timeline"></slot>
    </div>

    <div class="graph-container">
      <slot name="graph"></slot>
    </div>
  </div>
`;

export class VsCalendarDay extends HTMLElement {
  /** @type {HTMLElement} main container of this custom element */
  #container;

  /** @type {number} */
  #totalMinutes;

  constructor() {
    super();

    this.attachShadow({ mode: "open" }).append(
      document.importNode(template.content, true)
    );

    this.#container =
      this.shadowRoot?.querySelector(".main-container") ?? template;

    this.#totalMinutes = 1440;
  }

  connectedCallback() {
    const slots = this.shadowRoot?.querySelectorAll("slot");

    slots?.forEach(this.#setParentToChildren);
  }

  /** @param {Element} slot */
  #setParentToChildren = (slot) => {
    if (!(slot instanceof HTMLSlotElement)) return;

    slot.addEventListener("slotchange", (evt) => {
      slot.assignedElements().forEach((element) => {
        if (!(element instanceof VsPeriod || element instanceof VsTimeline))
          return;

        element.setParent(this);
      });
    });
  };

  // get the number of pixels per minute
  get pixelStep() {
    return this.#container.scrollHeight / this.#totalMinutes;
  }

  get totalMinutes() {
    return this.#totalMinutes;
  }

  get height() {
    return this.#container.clientHeight;
  }

  get top() {
    return (
      this.#container.getBoundingClientRect().top +
      document.documentElement.scrollTop
    );
  }
}

customElements.define("vs-calendar-day", VsCalendarDay);
