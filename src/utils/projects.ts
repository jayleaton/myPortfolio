import { IconType } from "react-icons"
import { SiNextdotjs, SiTypescript, SiKotlin, SiFirebase, SiExpress, SiRubyonrails } from "react-icons/si"
import { TbBrandSupabase } from "react-icons/tb"
import { FaAws, FaStripe, FaReact, FaDocker } from "react-icons/fa"
import { IoLogoJavascript } from "react-icons/io"

export interface Project {
  website: string
  github: string | null
  name: string
  image: string
  about: string
  stack: IconType[]
}

export const projects: Project[] = [
  {
    website: "https://simplifiedpdf.com",
    github: null,
    name: "<SimplifiedPDF />",
    image: "/simplifiedpdf.png",
    about:
      "SimplifiedPDF is a browser-based PDF editor for adding text, filling forms, and signing documents without signups, subscriptions, or server uploads. It is built to keep files on-device while still feeling fast and polished.",
    stack: [SiNextdotjs, SiTypescript],
  },
  {
    website: "https://simplifiedqr.com",
    github: null,
    name: "<SimplifiedQR />",
    image: "/simplifiedqr.png",
    about:
      "SimplifiedQR is a private QR generator for URLs, WiFi, contact cards, SMS, email, and plain text with free PNG and SVG exports. It focuses on a simple flow, fast generation, and keeping user data in the browser.",
    stack: [SiNextdotjs, SiTypescript],
  },
  {
    website: "https://doriracers.com",
    github: null,
    name: "<Doriracers />",
    image: "/dori.png",
    about:
      "Doriracers is the social platform for automotive enthusiasts to find and share their favorite car races and events. It is built on Next.js and the new app router along with supabase and postgres.",
    stack: [TbBrandSupabase, SiNextdotjs, SiTypescript, FaAws, FaDocker],
  },
  {
    website: "https://pricemates.net",
    github: null,
    name: "<Pricemates />",
    image: "/pricemates.png",
    about:
      "Pricemates allows users to save money on everyday purchases by using price alerts and viewing historical" +
      " price data. Pricemates is built on Next.js and the new app router along with supabase and postgres. We also" +
      " integrated stripe and observability so I am completely hands off now to focus on other projects.",
    stack: [TbBrandSupabase, SiNextdotjs, SiTypescript, FaStripe, FaAws, FaDocker],
  },
  {
    website: "https://globalpoker.com",
    github: null,
    name: "<Global Poker />",
    image: "/global-poker.png",
    about:
      `Global Poker is an industry leading online casino and is a part of the VGW line of offers, Working with` +
      ` the Poker team allowed me to work on several production projects including a significant` +
      ` contribution to an in house game "Ring Master"`,
    stack: [SiNextdotjs, SiTypescript, SiKotlin, FaAws, FaDocker],
  },
  {
    website: "https://noteteams.netlify.app/",
    github: "https://github.com/oldbettie/notepad",
    name: "<Team Notes />",
    image: "/teamnotes.png",
    about:
      "Team Notes was the final project of my software engineering bootcamp, team notes is a simplified miro clone." +
      "  Built using React for the front-end and Express/mySQL back-end that can also be found on my github. Best" +
      " experienced on desktop due to the canvas style of the project.",
    stack: [SiExpress, FaReact, IoLogoJavascript],
  },
  {
    website: "https://project-2-5825e.web.app/",
    github: "https://github.com/oldbettie/Project-2",
    name: "<Petpals />",
    image: "/petpals.png",
    about:
      "Bootcamp project designed as a pet adoption clone of Tinder where u can swipe left or right on pets then" +
      " communicate with the owner about adoption. built using React, FireBase, Framer Motion, Javascript and " +
      "SCSS Modules with a mobile first design.",
    stack: [SiFirebase, FaReact, IoLogoJavascript],
  },
  {
    website: "https://jayleaton-cookbook.herokuapp.com/",
    github: "https://github.com/oldbettie/cookbook",
    name: "<Cookbook />",
    image: "/cookbook.png",
    about:
      "Recipe finder app to find, post, edit or delete recipes utilizing RESTful best practices to build a CRUD app using Ruby on Rails, Heroku and Cloudinary ",
    stack: [SiRubyonrails, IoLogoJavascript],
  },
  {
    website: "https://oldbettie.github.io/tictactoe/",
    github: "https://github.com/oldbettie/tictactoe",
    name: "<Tictactoe />",
    image: "/tictactoedark.png",
    about: "Simple tictactoe game using basic JS and jQuery with added light/dark mode button. ",
    stack: [IoLogoJavascript],
  },
]
