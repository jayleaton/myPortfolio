import { Navigation } from "@/components/Navigation"
import { Content } from "@/components/Content"
import { BannerSection } from "@/components/BannerSection"
import { Companies } from "@/components/Companies"
import { Projects } from "@/components/Projects"
import { Skills } from "@/components/Skills"
import { About } from "@/components/About"
import { Footer } from "@/components/Footer"

export default function Home() {
  return (
    <div className="relative overflow-x-clip">
      <div className="top-0 w-full">
        <Content>
          <Navigation />
          <BannerSection />
        </Content>
      </div>
      <div className="h-full w-full pb-10">
        <Content>
          <Companies />
          <Projects />
          <Skills />
          <About />
          <Footer />
        </Content>
      </div>
    </div>
  )
}
