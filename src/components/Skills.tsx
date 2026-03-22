import { FaGithub, FaStripe, FaDocker } from "react-icons/fa"
import { DiCss3, DiHtml5, DiReact, DiRuby, DiSass, DiNodejs, DiJava } from "react-icons/di"
import {
  SiNextdotjs,
  SiExpress,
  SiPostgresql,
  SiFirebase,
  SiKotlin,
  SiTerraform,
  SiObservable,
  SiTypescript,
  SiJavascript,
  SiTailwindcss,
  SiAmazonaws,
  SiJquery,
  SiPixiv,
  SiThreedotjs,
  SiUnity,
  SiSupabase,
  SiAuth0,
  SiJira,
  SiCircleci,
} from "react-icons/si"
import { TbBrandSupabase } from "react-icons/tb"
export function Skills() {
  return (
    <section id="Skills" className="section-shell relative scroll-mt-28 py-20">
      <div className="section-grid absolute inset-0 opacity-60" />
      <div className="section-panel relative overflow-hidden rounded-[32px] border border-white/10 px-6 py-10 md:px-8 xl:px-10">
        <div className="mb-12 flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
          <div className="max-w-2xl">
            <p className="section-kicker">Tools I reach for</p>
            <p className="mt-4 leading-relaxed text-font_main">
              Product-focused full-stack work across frontend, backend, cloud, data, and delivery tooling.
            </p>
          </div>
        </div>
        <div className="grid grid-cols-3 gap-4 sm:grid-cols-4 lg:grid-cols-6">
          {[
          TbBrandSupabase,
          FaGithub,
          SiJavascript,
          FaDocker,
          FaStripe,
          DiCss3,
          DiHtml5,
          DiReact,
          SiNextdotjs,
          DiRuby,
          SiJquery,
          SiTypescript,
          DiNodejs,
          SiTailwindcss,
          SiAmazonaws,
          SiExpress,
          DiSass,
          SiPostgresql,
          SiFirebase,
          SiKotlin,
          SiTerraform,
          SiObservable,
          DiJava,
          SiPixiv,
          SiThreedotjs,
          SiUnity,
          SiSupabase,
          SiAuth0,
          SiJira,
          SiCircleci,
        ].map((Icon, index) => (
          <div
            key={index}
            className="group flex min-h-[110px] items-center justify-center rounded-[24px] border border-white/10 bg-slate-950/30 text-[40px] text-font_main shadow-[0_20px_50px_rgba(2,6,23,0.18)] transition-all duration-300 hover:-translate-y-1 hover:border-font_highlight/30 hover:text-font_highlight md:text-[52px]"
          >
            <Icon className="transition-transform duration-300 group-hover:scale-110" />
          </div>
        ))}
        </div>
      </div>
    </section>
  )
}
