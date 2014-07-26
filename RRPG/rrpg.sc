package console := import("console");

type Stage := struct {
 name:string,
 mobs:int,
 maxMobHealth:int,
 next:Stage
}

type Player := struct {
 name:string,
 hp:int,
 stage:Stage
}

func ProcessFight(player:Player, mob:Player) {
}

func Process(player:Player) {
 if player->stage = nil then return;
 var mobs := player->stage->mobs;

 console.Log("Player entered stage " $ player->stage->name $ "\n");

 while player->stage->mobs > 0 do {
  player->stage->mobs := mobs;
  //player->stage->mobs := player->stage->mobs - 1;
 }

 player->stage := player->stage->next;
 Process(player);
}

func main() {

 var stage2 := Stage { "Stage 2", 10, 40, nil };
 var stage1 := Stage { "Stage 1", 5, 20, stage2 };

 var blake := Player{"Blake", 100, stage1 };

 console.Log("Randomo the game-yo\n");
 Process(blake);
}
