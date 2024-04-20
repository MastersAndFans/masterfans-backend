package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/MastersAndFans/masterfans-backend/pkg/helpers"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type AuctionHandler struct {
	AuctionRepo repository.IAuctionRepository
	UserRepo repository.IUserRepository
}

func NewAuctionHandler(auctionRepo repository.IAuctionRepository, userRepo repository.IUserRepository) *AuctionHandler {
	return &AuctionHandler{AuctionRepo: auctionRepo, UserRepo: userRepo}
}


type UpdateBidderPayload struct {
	Auction_ID int64
	User_ID int64
	Bid int64
}

func (handler *AuctionHandler) ListAuctions(w http.ResponseWriter, r *http.Request) {
	auctions, err := handler.AuctionRepo.List(context.Background())

	auctions_json, err := json.Marshal(auctions)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create JSON"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(auctions_json)
}

func (handler *AuctionHandler) GetAuctionById(w http.ResponseWriter, r *http.Request) {
	id_string := chi.URLParam(r, "id")

	// convert id to int64
	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID must contain only digits"})
		return
	}

	auction, err := handler.AuctionRepo.FindById(context.Background(), id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Auction with ID %d not found", id)})
		return
	}

	auction_json, err := json.Marshal(auction)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create JSON"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(auction_json)
}

func (handler *AuctionHandler) CreateAuction(w http.ResponseWriter, r *http.Request) {
	var auction models.Auction
	if err := json.NewDecoder(r.Body).Decode(&auction); err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, err.Error())
		return
	}

	if auction.StartDate.After(auction.EndDate) {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Start date cannot be after end date")
		return
	}

	err := handler.AuctionRepo.Create(r.Context(), &auction)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create auction")
		return
	}

	response := map[string]string{
		"message": "Auction created successfully",
		"auction_id": fmt.Sprint(auction.ID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create JSON")
	}
}

func (handler *AuctionHandler) UpdateBidder(w http.ResponseWriter, r *http.Request) {
	var request UpdateBidderPayload
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, err.Error())
		return
	}

	auction, err := handler.AuctionRepo.FindById(r.Context(), request.Auction_ID)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, err.Error())
		return
	}

	// check if bid is lower
	if auction.StartingPrice <= request.Bid {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Proposed bid is not lower than current bid")
		return
	}

	user, err := handler.UserRepo.FindById(r.Context(), request.User_ID)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, err.Error())
		return
	}


	auction.StartingPrice = request.Bid
	auction.Participants = append(auction.Participants, user)

	err = handler.AuctionRepo.Update(r.Context(), auction)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(auction); err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create JSON")
	}
}

func (handler *AuctionHandler) DeleteAuction(w http.ResponseWriter, r *http.Request) {
	id_string := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID must contain only digits"})
		return
	}

	err = handler.AuctionRepo.Delete(r.Context(), id)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{
		"message": "Auction deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println(err)
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create JSON")
	}
}
