package controller

import (
	"errors"
	"general_game/gcontroller"
	"general_game/gmodel"
	"general_game/gutils"
	"log"

	"gopkg.in/mgo.v2/bson"
)

// SetStart t
func SetStart(chatID int64) error {
	//database connection
	database, session, err := gcontroller.GetDB("telegram", gutils.DBTelegramURITeen)

	defer session.Close()

	if err != nil {
		log.Print("to get database, in NewPlayer")
		return errors.New("database")
	}

	subCol := database.C(gutils.SUBSCRIBERS)

	var result gmodel.TelegramSubscriber
	err = subCol.Find(bson.M{"chatID": chatID}).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			err = subCol.Insert(bson.M{"chatID": chatID, "status": gutils.StartLsofStatus})
			if err != nil {
				log.Print(err.Error())
				return errors.New("not inserted")
			}
			return nil
		}
	}

	if result.Status == gutils.StopStatus {
		err = subCol.Update(bson.M{}, bson.M{"$set": bson.M{"chatID": chatID, "status": gutils.StartLsofStatus}})
		if err != nil {
			return errors.New("not inserted")
		}
	}

	return err
}

// SetStop t
func SetStop(chatID int64) error {
	//database connection
	database, session, err := gcontroller.GetDB("telegram", gutils.DBTelegramURITeen)

	defer session.Close()

	if err != nil {
		log.Print("to get database, in NewPlayer")
		return errors.New("database")
	}

	subCol := database.C(gutils.SUBSCRIBERS)

	err = subCol.Remove(bson.M{"chatID": chatID})
	if err != nil {
		return errors.New("not inserted")
	}

	return err
}

// SetLsofStop t
func SetLsofStop(chatID int64) error {
	//database connection
	database, session, err := gcontroller.GetDB("telegram", gutils.DBTelegramURITeen)

	defer session.Close()

	if err != nil {
		log.Print("to get database, in NewPlayer")
		return errors.New("database")
	}

	subCol := database.C(gutils.SUBSCRIBERS)

	var result gmodel.TelegramSubscriber
	err = subCol.Find(bson.M{"chatID": chatID}).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			err = subCol.Insert(bson.M{"chatID": chatID, "status": gutils.StopLsofStatus})
			if err != nil {
				return errors.New("not inserted")
			}
			return nil
		}
		return err
	}

	if result.Status != gutils.StopLsofStatus {
		err = subCol.Update(bson.M{}, bson.M{"$set": bson.M{"chatID": chatID, "status": gutils.StopLsofStatus}})
		if err != nil {
			return errors.New("not inserted")
		}
	}

	return err
}

// GetAllLsof used to get all subs lsof
func GetAllLsof() ([]gmodel.TelegramSubscriber, error) {
	//database connection
	database, session, err := gcontroller.GetDB("telegram", gutils.DBTelegramURITeen)

	defer session.Close()

	if err != nil {
		log.Print("to get database, in NewPlayer")
		return nil, errors.New("database")
	}

	subCol := database.C(gutils.SUBSCRIBERS)

	var result []gmodel.TelegramSubscriber
	err = subCol.Find(bson.M{"status": gutils.StartLsofStatus}).All(&result)

	return result, nil
}

// GetAllActive used to get all active subs
func GetAllActive() ([]gmodel.TelegramSubscriber, error) {
	//database connection
	database, session, err := gcontroller.GetDB("telegram", gutils.DBTelegramURITeen)

	defer session.Close()

	if err != nil {
		log.Print("to get database, in NewPlayer")
		return nil, errors.New("database")
	}

	subCol := database.C(gutils.SUBSCRIBERS)

	var result []gmodel.TelegramSubscriber
	// "status": gutils.StartLsofStatus
	err = subCol.Find(bson.M{"status": bson.M{"$ne": gutils.StopStatus}}).All(&result)
	log.Print(len(result))
	return result, nil
}
