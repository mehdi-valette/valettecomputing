/** @typedef {import ("./calendar-day.js").VsCalendarDay} VsCalendarDay */

import { minutesToHours } from "../minutes-to-hours.js";

const containerTemplate = document.createElement("template");
containerTemplate.innerHTML = `
  <style>
    .container {
      position: relative;
      width: 10rem;
    }

    .line {
      border-top: 1px solid blue;
      position: absolute;
      width: 100%;
    }
  </style>

  <div class="container"></div>
`;

const lineTemplate = document.createElement("template");
lineTemplate.innerHTML = `
  <div class="line"></div>
`;

export class VsTimeline extends HTMLElement {
  /** @type {VsCalendarDay | null} parent */
  #parent;

  /** @type {HTMLElement} the container of the timeline */
  #container;

  constructor() {
    super();

    this.#parent = null;

    this.attachShadow({ mode: "open" }).appendChild(
      document.importNode(containerTemplate.content, true)
    );

    this.#container =
      this.shadowRoot?.querySelector(".container") ?? containerTemplate;
  }

  /** @param {VsCalendarDay} parent */
  setParent = (parent) => {
    this.#parent = parent;

    this.#createTimeline();
  };

  #createTimeline = () => {
    const pixelStep = this.#parent?.pixelStep ?? 0;

    for (let i = 0; i < 48; i++) {
      const line = document.importNode(lineTemplate.content, true);

      const div = line.querySelector("div");
      console.log(div);

      if (div == null) return;

      div.style.top = pixelStep * 30 * i + "px";
      div.append(document.createTextNode(minutesToHours(i * 30)));

      this.#container.appendChild(line);
    }
  };
}

customElements.define("vs-timeline", VsTimeline);
