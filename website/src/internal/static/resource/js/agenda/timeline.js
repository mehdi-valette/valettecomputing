/** @typedef {import ("./calendar-day.js").VsCalendarDay} VsCalendarDay */

import { minutesToHours } from "../minutes-to-hours.js";

const containerTemplate = document.createElement("template");
containerTemplate.innerHTML = `
  <style>
    .container {
      position: relative;
      width: 5rem;
    }

    .line {
      border-top: 1px solid blue;
      position: absolute;
      width: 100%;
    }

    .time {
      margin-top: -0.6em;
      background-color: white;
      width: 3rem;
    }
  </style>

  <div class="container"></div>
`;

const lineTemplate = document.createElement("template");
lineTemplate.innerHTML = `
  <div class="line">
    <div class="time"></div>
  </div>
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

    for (let i = 0; i <= 1440; i += 60) {
      const template = document.importNode(lineTemplate.content, true);

      const line = template.querySelector(".line");
      const time = template.querySelector(".time");

      if (!(line instanceof HTMLElement && time instanceof HTMLElement)) return;

      line.style.top = pixelStep * i + "px";
      time.append(document.createTextNode(minutesToHours(i)));

      this.#container.appendChild(template);
    }
  };
}

customElements.define("vs-timeline", VsTimeline);
