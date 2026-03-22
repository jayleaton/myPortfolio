import { Button } from "@/components/Button"
import { SiNextdotjs, SiPostgresql, SiTypescript } from "react-icons/si"

const stats = [
  { value: "15", label: "Years in business" },
  { value: "30+", label: "Projects completed" },
  { value: "5", label: "Years shipping software" },
  { value: "12+", label: "Companies helped" },
]

export function BannerSection() {
  return (
    <section className="hero-shell relative py-12 md:py-16 xl:py-20">
      <div className="hero-grid absolute inset-0 opacity-70" />
      <div className="hero-orb hero-orb-left" />
      <div className="hero-orb hero-orb-right" />

      <div className="relative z-10 grid items-center gap-14 xl:grid-cols-[1.05fr_0.95fr]">
        <div className="hero-panel reveal-up flex flex-col gap-6 rounded-[32px] border border-white/10 bg-white/[0.03] p-6 shadow-[0_40px_120px_rgba(15,23,42,0.55)] backdrop-blur-xl md:p-8 xl:p-10">
          <div className="reveal-up delay-1 flex flex-wrap items-center gap-3">
            <span className="font-code text-xs uppercase tracking-[0.35em] text-font_highlight md:text-sm">
              Software Engineer
            </span>
            <span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 font-code text-xs text-font_main">
              SaaS, AI, product delivery
            </span>
          </div>

          <div className="reveal-up delay-2 flex flex-col gap-4">
            <p className="font-code text-xs uppercase tracking-[0.45em] text-font_dark md:text-sm">
              Building software with business context
            </p>
            <h1 className="max-w-3xl font-code text-[2.6rem] font-bold leading-[0.95] text-white md:text-[4.6rem] xl:text-[5.2rem]">
              Code.
              <br />
              Create.
              <br />
              <span className="bg-[linear-gradient(135deg,#67e8f9_0%,#2dd4bf_35%,#60a5fa_100%)] bg-clip-text text-transparent">
                Innovate. Deliver.
              </span>
            </h1>
          </div>

          <div className="reveal-up delay-3 grid gap-4 text-[17px] leading-8 text-font_main md:max-w-2xl md:text-[18px]">
            <p>
              With 15 years in business management and 5 years shipping software, I focus on pragmatic systems that
              move fast without becoming fragile. SaaS and micro SaaS are where that mindset compounds.
            </p>
            <p>
              I use modern AI tooling to accelerate delivery, but the bar stays the same: maintainable architecture,
              clear product thinking, and code a team can keep extending long after launch.
            </p>
          </div>

          <div className="reveal-up delay-4 flex flex-col gap-4 pt-2 sm:flex-row sm:items-center">
            <Button
              className="rounded-2xl border border-font_highlight/40 bg-font_highlight px-6 py-3 font-code text-sm font-bold tracking-[0.2em] text-background transition-all duration-300 hover:-translate-y-1 hover:shadow-[0_18px_45px_rgba(45,212,191,0.35)]"
              path="mailto: jayleaton@gmail.com"
            >
              LET&apos;S CHAT
            </Button>
            <p className="font-code text-sm text-font_dark">
              Available for product-minded engineering roles and freelance builds.
            </p>
          </div>

          <div className="reveal-up delay-5 grid gap-4 pt-4 sm:grid-cols-2 xl:grid-cols-4">
            {stats.map((stat) => (
              <div
                key={stat.label}
                className="rounded-2xl border border-white/8 bg-slate-950/30 px-4 py-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.04)]"
              >
                <p className="font-interface text-3xl font-semibold text-font_header md:text-4xl">{stat.value}</p>
                <p className="mt-2 max-w-[12ch] font-code text-xs uppercase tracking-[0.2em] text-font_main">
                  {stat.label}
                </p>
              </div>
            ))}
          </div>
        </div>

        <div className="relative z-10 flex min-h-[520px] items-center justify-center xl:min-h-[640px]">
          <div className="hero-visual reveal-fade delay-2 relative aspect-square w-full max-w-[620px] rounded-full border border-white/10 bg-[radial-gradient(circle_at_center,rgba(148,163,184,0.08),rgba(15,23,42,0.02)_60%,transparent_70%)]">
            <div className="absolute inset-[12%] rounded-full border border-white/5 bg-white/[0.02]" />
            <div className="absolute inset-[21%] rounded-full border border-white/10 bg-background_highlight/35 backdrop-blur-md" />
            <img
              src="/code-bubbles.png"
              alt="Text bubbles floating in the hero artwork"
              className="hero-image absolute left-1/2 top-1/2 w-[88%] max-w-[540px] -translate-x-1/2 -translate-y-1/2 drop-shadow-[0_35px_60px_rgba(15,23,42,0.75)]"
            />
          </div>

          <div className="reveal-fade delay-3 absolute left-0 top-[58%]">
            <div className="float-badge float-badge-left">
              <div className="hero-badge hero-badge-next">
                <SiNextdotjs />
              </div>
            </div>
            <p className="mt-3 font-code text-xs uppercase tracking-[0.3em] text-font_main">Next.js</p>
          </div>

          <div className="reveal-fade delay-4 absolute right-[2%] top-[10%]">
            <div className="float-badge float-badge-top">
              <div className="hero-badge hero-badge-postgresql">
                <SiPostgresql />
              </div>
            </div>
            <p className="mt-3 font-code text-xs uppercase tracking-[0.3em] text-font_main">PostgreSQL</p>
          </div>

          <div className="reveal-fade delay-5 absolute bottom-[6%] right-[6%]">
            <div className="float-badge float-badge-right">
              <div className="hero-badge hero-badge-typescript">
                <SiTypescript />
              </div>
            </div>
            <p className="mt-3 font-code text-xs uppercase tracking-[0.3em] text-font_main">TypeScript</p>
          </div>
        </div>
      </div>
    </section>
  )
}
