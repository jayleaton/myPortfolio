export type YouTubeVideo = {
  id: string
  title: string
}

const YOUTUBE_CHANNEL_ID = "UC-AGtYCONF0Yu0uK_fL1j1Q"
const YOUTUBE_FEED_URL = `https://www.youtube.com/feeds/videos.xml?channel_id=${YOUTUBE_CHANNEL_ID}`

const FALLBACK_VIDEOS: YouTubeVideo[] = [
  {
    id: "dHiL7HNphJc",
    title: "The Best AI Agent Workflow Right NOW!",
  },
  {
    id: "N-aYRmXddbw",
    title: "Code is Only Part of the Problem (Learn This Now!)",
  },
  {
    id: "aOPcq-u8Siw",
    title: "Your Coding Skills Are Obsolete | Here's Why",
  },
  {
    id: "TuNoSyWRpXI",
    title: "SCARY Bugs Hidden in Popular Next.js Repos",
  },
]

function decodeXmlEntities(value: string) {
  return value
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
}

function parseFeedEntries(xml: string) {
  const entryMatches = xml.match(/<entry>[\s\S]*?<\/entry>/g) ?? []

  return entryMatches
    .map((entry) => {
      const id = entry.match(/<yt:videoId>([^<]+)<\/yt:videoId>/)?.[1]?.trim()
      const title = entry.match(/<title>([\s\S]*?)<\/title>/)?.[1]?.trim()

      if (!id || !title) {
        return null
      }

      return {
        id,
        title: decodeXmlEntities(title),
      }
    })
    .filter((video): video is YouTubeVideo => video !== null)
}

export async function getLatestYouTubeVideos(limit = 4): Promise<YouTubeVideo[]> {
  try {
    const response = await fetch(YOUTUBE_FEED_URL, {
      cache: "no-store",
    })

    if (!response.ok) {
      throw new Error(`YouTube feed request failed with ${response.status}`)
    }

    const xml = await response.text()
    const videos = parseFeedEntries(xml).slice(0, limit)

    if (videos.length > 0) {
      return videos
    }
  } catch (error) {
    console.error("Failed to load latest YouTube videos", error)
  }

  return FALLBACK_VIDEOS.slice(0, limit)
}
