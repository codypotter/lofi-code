---
title: Being right is overrated
slug: being-right-is-overrated
summary: >-
  It's a high-stakes discussion in a boardroom at BigTek. The
    stakeholders are not happy with your team's progress, because they needed the
    "Real-Time Activity Feed" ™️ like two weeks ago. Of course, your team has been
    focused on that big, surprise refactor for the last few weeks, but now that
    that Archimedes and FooBar are PEEPEEPOOPOO compliant, we're good to start
    that new, all-important feature.
date: 2024-11-03T19:23:00.000Z
headerImage: https://loficode.com/media/being-right-is-overrated-3-1.png
openGraphImage: https://loficode.com/media/being-right-is-overrated-16-9.png
tags:
  - soft-skills
  - programming
---
It's a high-stakes discussion in a boardroom at BigTek. The stakeholders are not happy with your team's progress, because they needed the "Real-Time Activity Feed" ™️ like two weeks ago to meet a major marketing deadline. Of course, your team has been focused on that big, surprise refactor for the last few weeks. But now that Archimedes and FooBar are PEEPEEPOOPOO compliant, we're good to start that new, all-important feature.

Now it's your time to shine. You've been thinking deeply and considering every single edge case during your morning shower.

"I think this calls for a MariaDB with a centralized 'Activity' table. Every user action will be stored there. We can query that table in real time to retrieve the latest actions and display them to users."

Your buddy Taylor chimes in, "We should call it 'activities' instead because singular database names are for murderers and psychopaths."

Cameron, who is often quiet in meetings raises their hand. Oh no, you think, Not Cameron, they always prove me wrong.

"I don't think that will work, at least not without some modifications."

Your heart sinks. You have fucked up. You close your eyes for a brief moment, anticipating that everyone in the room is now turning to you and laughing. They all hate you. You absolute buffoon, you forgot to bring your clown shoes to work today. How did you manage to fuck up this bad. How could this have happened? You considered every scenario. You thought of all the foreign keys. You even drew a UML diagram in the steam on your bathroom mirror.

Cameron calmly states, "This would be a very write-heavy workload. Putting everything in a table would cause a serious bottleneck at our scale. I think we can make the activities table work though. If we store all recent activity in a Redis, we can use some horizontal scaling for fast reads and use an SQS queue to add eventual persistence to the database, I think that schema should work, but we also get better speed and scalability."

You open your eyes, expecting your whole team to be staring at you with disappointment. The pit in your stomach disappears, because no one is looking your way. Everyone is looking at Cameron. Reese stands up and starts diagramming Cameron's solution on the whiteboard.

How could this be?

---

Being right is overrated. Getting it right matters more.

There are no high scores within a team, and there are no "better winners" on a team. Teams win and lose together. It's a better strategy to surround yourself with people who are right a lot than to bear the burden never getting it wrong.

The president has a cabinet, CEOs have boards, comprehensive care patients have teams of surgeons. You have Cameron, Taylor, Reese, and Bailey.

I think this is an example where losing is winning. Sure, you were wrong champ, but at the end of the day your team got it right. There is a part of you that is actually grateful that Cameron spoke up.

Wouldn't it have been much worse if your team had to spend another two weeks backtracking on your initial design after making a delivery date promise to your stakeholders?

The benefits of being wrong go even further than that:

* You're going to learn from your mistakes a lot.
* You're going to learn how to handle setbacks pretty good.
* You're going to strengthen those muscles of humility.
* Admitting you're wrong demonstrates your integrity and self-awareness, which builds credibility and trustworthiness.

Embrace that uncomfortable feeling of being wrong. Admit when you're wrong so your team can always be right.
