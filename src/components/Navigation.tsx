"use client"

import React, { useEffect, useState } from "react"
import { Button } from "@/components/Button"
import { UseMediaQuery } from "@/hooks/useMediaQuery"
import { FaGithub, FaLinkedin } from "react-icons/fa"
import { BiLogoUpwork } from "react-icons/bi"
import { HamburgerMenu } from "@/components/HamburgerMenu"

const getIsDesktop = (): boolean => {
  return UseMediaQuery(900)
}

enum NavSections {
  "COMPANIES" = "Companies",
  "PROJECTS" = "Projects",
  "SKILLS" = "Skills",
  "ABOUT" = "About",
  "RESUME" = "Resume",
}

export function Navigation() {
  const isDesktop = getIsDesktop()
  const [activeSection, setActiveSection] = useState<NavSections>(NavSections.COMPANIES)

  useEffect(() => {
    const sectionMap: Record<string, NavSections> = {
      Companies: NavSections.COMPANIES,
      Projects: NavSections.PROJECTS,
      Skills: NavSections.SKILLS,
      About: NavSections.ABOUT,
    }

    const sectionObserver = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const nextSection = sectionMap[entry.target.id]
            if (nextSection) setActiveSection(nextSection)
          }
        })
      },
      { threshold: 0.5 }
    )

    const sections = document.querySelectorAll("#Companies, #Projects, #Skills, #About")
    sections.forEach((section) => {
      sectionObserver.observe(section)
    })

    return () => {
      sectionObserver.disconnect()
    }
  }, [])

  if (isDesktop) {
    return (
      <div className="flex w-full justify-between items-center py-6">
        <h1 className="font-interface text-3xl font-bold text-white tracking-tight">
          JJ <span className="text-font_highlight">Eaton</span>
        </h1>
        <div className="w-fit">
          <div className="text-right lg:space-x-8 space-x-4 pl-5 text-[16px] flex items-center">
            <Button
              path="#Companies"
              className={`font-code transition-colors hover:text-font_highlight ${
                activeSection === NavSections.COMPANIES ? "text-font_highlight" : "text-font_light"
              }`}
              click={() => setActiveSection(NavSections.COMPANIES)}
            >
              {activeSection === NavSections.COMPANIES ? `< Companies />` : `Companies`}
            </Button>
            <Button
              path="#Projects"
              className={`font-code transition-colors hover:text-font_highlight ${
                activeSection === NavSections.PROJECTS ? "text-font_highlight" : "text-font_light"
              }`}
              click={() => setActiveSection(NavSections.PROJECTS)}
            >
              {activeSection === NavSections.PROJECTS ? `< Projects />` : `Projects`}
            </Button>
            <Button
              path="#Skills"
              className={`font-code transition-colors hover:text-font_highlight ${
                activeSection === NavSections.SKILLS ? "text-font_highlight" : "text-font_light"
              }`}
              click={() => setActiveSection(NavSections.SKILLS)}
            >
              {activeSection === NavSections.SKILLS ? `< Skills />` : `Skills`}
            </Button>
            <Button
              path="#About"
              className={`font-code transition-colors hover:text-font_highlight ${
                activeSection === NavSections.ABOUT ? "text-font_highlight" : "text-font_light"
              }`}
              click={() => setActiveSection(NavSections.ABOUT)}
            >
              {activeSection === NavSections.ABOUT ? `< About />` : `About`}
            </Button>
            <Button
              path="resume.pdf"
              target="_blank"
              className={`font-code transition-colors hover:text-font_highlight ${
                activeSection === NavSections.RESUME ? "text-font_highlight" : "text-font_light"
              }`}
              click={() => setActiveSection(NavSections.RESUME)}
            >
              {activeSection === NavSections.RESUME ? `< Resume />` : `Resume`}
            </Button>
          </div>
        </div>

        <div className="flex gap-4">
          <Button path="https://github.com/oldbettie" target="_blank" className="hover:text-white text-font_light transition-colors">
            <div className="flex justify-center items-center gap-2">
              <FaGithub className="text-xl" />
              <span className="hidden lg:inline">Github</span>
            </div>
          </Button>
          <Button path="https://www.linkedin.com/in/jayj-eaton" target="_blank" className="hover:text-blue-400 text-font_light transition-colors">
            <div className="flex justify-center items-center gap-2">
              <FaLinkedin className="text-xl" />
              <span className="hidden lg:inline">Linkedin</span>
            </div>
          </Button>
          <Button path="https://www.upwork.com/freelancers/~01bb2932f2d3e396a6" target="_blank" className="hover:text-green-400 text-font_light transition-colors">
            <div className="flex justify-center items-center gap-2">
              <BiLogoUpwork className="text-2xl" />
              <span className="hidden lg:inline">Upwork</span>
            </div>
          </Button>
        </div>
      </div>
    )
  } else {
    return (
      <div className="flex w-full justify-between pt-5">
        <h1 className="w-full text-2xl text-font_main font-interface">JJ Eaton</h1>
        <HamburgerMenu />
      </div>
    )
  }
}
