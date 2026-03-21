import { Button } from "@/components/Button"
import Link from "next/link"
import { FaBuilding, FaGlobe, FaYoutube } from "react-icons/fa"

const latestVideos = [
  {
    id: "N-aYRmXddbw",
    title: "Code is Only Part of the Problem (Learn This Now!)",
  },
  {
    id: "aOPcq-u8Siw",
    title: "SCARY Bugs Hidden in Popular Next.js Repos",
  },
  {
    id: "TuNoSyWRpXI",
    title: "Don't Waste Months on SaaS Ideas",
  },
]

const companies = [
  {
    name: "Build Things YouTube",
    href: "https://www.youtube.com/@build-things",
    icon: FaYoutube,
    cta: "Watch",
    ctaIcon: FaYoutube,
    accent: "text-red-400",
    border: "border-red-400/20",
    description:
      "A place where I share how I build software, think through product decisions, and turn ideas into working systems.",
    videos: latestVideos,
  },
  {
    name: "Buildthings.co",
    href: "https://buildthings.co/",
    icon: FaBuilding,
    cta: "Website",
    ctaIcon: FaGlobe,
    accent: "text-font_highlight",
    border: "border-font_highlight/20",
    description:
      "My software agency focused on designing and shipping reliable products for founders and teams that need execution without unnecessary overhead.",
  },
]

export function Companies() {
  return (
    <div id="Companies" className="py-20 relative scroll-mt-24">
      <div className="flex flex-col md:flex-row justify-between items-end mb-12 gap-6">
        <div className="md:w-1/2">
          <p className="text-font_main max-w-xl leading-relaxed">
            More places to watch what I am building, see what I am working on, and get in touch when you want to collaborate.
          </p>
        </div>
        <h2 className="text-[35px] md:text-[50px] font-bold font-code text-right leading-tight">
          Find Me <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-font_highlight to-blue-500">
            Elsewhere
          </span>
        </h2>
      </div>

      <div className="flex flex-col gap-6">
        {companies.map((company) => {
          const Icon = company.icon
          const CTAIcon = company.ctaIcon
          return (
            <div
              key={company.name}
              className={`group relative overflow-hidden rounded-2xl border ${company.border} bg-background_highlight/40 backdrop-blur-sm p-6 md:p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl hover:shadow-font_highlight/10`}
            >
              <div className="absolute inset-y-0 left-0 w-1 bg-gradient-to-b from-transparent via-font_highlight/60 to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
              <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6">
                <div className="max-w-2xl">
                  <div className="flex items-center gap-3 mb-4">
                    <div className={`rounded-full border border-white/10 bg-white/5 p-3 ${company.accent}`}>
                      <Icon size={22} />
                    </div>
                    <h3 className="text-white text-2xl md:text-3xl font-code font-bold group-hover:text-font_highlight transition-colors">
                      {company.name}
                    </h3>
                  </div>
                  <p className="text-font_main text-base md:text-lg leading-relaxed">{company.description}</p>

                  {company.videos && (
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-6">
                      {company.videos.map((video) => (
                        <Link
                          key={video.id}
                          href={`https://www.youtube.com/watch?v=${video.id}`}
                          target="_blank"
                          className="group/video block overflow-hidden rounded-xl border border-white/10 bg-black/30 transition-all duration-300 hover:-translate-y-1 hover:border-red-400/40"
                        >
                          <div className="relative">
                            <img
                              src={`https://i.ytimg.com/vi/${video.id}/hqdefault.jpg`}
                              alt={video.title}
                              className="h-auto w-full object-cover transition-transform duration-500 group-hover/video:scale-105"
                            />
                            <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/10 to-transparent" />
                            <div className="absolute bottom-3 left-3 inline-flex items-center gap-2 rounded-full bg-black/70 px-3 py-1 text-xs font-code text-white">
                              <FaYoutube className="text-red-400" />
                              Latest video
                            </div>
                          </div>
                          <div className="p-4">
                            <p className="line-clamp-2 text-sm leading-relaxed text-font_light">{video.title}</p>
                          </div>
                        </Link>
                      ))}
                    </div>
                  )}
                </div>

                <Button
                  path={company.href}
                  target="_blank"
                  className="w-fit font-code text-base md:text-lg border border-white/10 text-font_light px-6 py-3 rounded-full hover:bg-white/5 hover:text-white transition-all"
                >
                  <span className="inline-flex items-center gap-3">
                    <CTAIcon />
                    {company.cta}
                  </span>
                </Button>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
