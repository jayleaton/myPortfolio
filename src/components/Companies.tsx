import { Button } from "@/components/Button"
import { getLatestYouTubeVideos } from "@/utils/youtube"
import Link from "next/link"
import { FaBuilding, FaGlobe, FaYoutube } from "react-icons/fa"

function getVideoVisibilityClass(index: number) {
  if (index === 0) {
    return "block"
  }

  if (index === 1) {
    return "hidden sm:block"
  }

  if (index === 2) {
    return "hidden lg:block"
  }

  return "hidden 2xl:block"
}

export async function Companies() {
  const latestVideos = await getLatestYouTubeVideos(4)
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

  return (
    <div id="Companies" className="py-20 relative scroll-mt-24">
      <div className="mb-12 flex flex-col items-center gap-6 text-center md:flex-row md:items-end md:justify-between md:text-left">
        <div className="md:w-1/2">
          <p className="mx-auto max-w-xl leading-relaxed text-font_main md:mx-0">
            More places to watch what I am building, see what I am working on, and get in touch when you want to collaborate.
          </p>
        </div>
        <h2 className="text-center text-[35px] font-bold font-code leading-tight md:text-right md:text-[50px]">
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
              <div className="flex flex-col gap-6 text-center md:flex-row md:items-center md:justify-between md:text-left">
                <div className="flex-1 min-w-0">
                  <div className="mb-4 flex items-center justify-center gap-3 md:justify-start">
                    <div className={`rounded-full border border-white/10 bg-white/5 p-3 ${company.accent}`}>
                      <Icon size={22} />
                    </div>
                    <h3 className="text-white text-2xl md:text-3xl font-code font-bold group-hover:text-font_highlight transition-colors">
                      {company.name}
                    </h3>
                  </div>
                  <p className="mx-auto max-w-2xl text-base leading-relaxed text-font_main md:mx-0 md:text-lg">
                    {company.description}
                  </p>

                  {company.videos && (
                    <div className="mt-6 flex flex-nowrap justify-center gap-4 overflow-hidden md:justify-start">
                      {company.videos.map((video, index) => (
                        <Link
                          key={video.id}
                          href={`https://www.youtube.com/watch?v=${video.id}`}
                          target="_blank"
                          className={`group/video w-full max-w-[210px] min-w-0 shrink-0 overflow-hidden rounded-xl border border-white/10 bg-black/30 transition-all duration-300 hover:-translate-y-1 hover:border-red-400/40 ${getVideoVisibilityClass(index)}`}
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
                  className="self-center w-fit rounded-full border border-white/10 px-6 py-3 font-code text-base text-font_light transition-all hover:bg-white/5 hover:text-white md:self-auto md:text-lg"
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
