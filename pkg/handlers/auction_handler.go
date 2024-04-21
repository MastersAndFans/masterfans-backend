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
	"time"
)

type AuctionHandler struct {
	AuctionRepo repository.IAuctionRepository
	UserRepo    repository.IUserRepository
}

func NewAuctionHandler(auctionRepo repository.IAuctionRepository, userRepo repository.IUserRepository) *AuctionHandler {
	return &AuctionHandler{AuctionRepo: auctionRepo, UserRepo: userRepo}
}

type CreateAuctionPayload struct {
	ProposerId    uint   `json:"proposer_id"`
	StartingPrice uint   `json:"starting_price"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Active        bool   `json:"active"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	City          string `json:"city"`
	Category      uint   `json:"category"`
}

type UpdateBidderPayload struct {
	Auction_ID int64
	User_ID    int64
	Bid        int64
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
	var payload CreateAuctionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, err.Error())
		return
	}

	const date_layout = "2006-01-02 15:04:05"
	start_date, err := time.Parse(date_layout, payload.StartDate)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Start date does not match yyyy-mm-dd hh:mm:ss format")
		return
	}

	end_date, err := time.Parse(date_layout, payload.EndDate)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, "End date does not match yyyy-mm-dd hh:mm:ss format")
		return
	}

	if start_date.After(end_date) {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Start date cannot be after end date")
		return
	}

	auction := models.Auction{
		ProposerID:    payload.ProposerId,
		Active:        payload.Active,
		StartingPrice: int64(payload.StartingPrice),
		StartDate:     start_date,
		EndDate:       end_date,
		Title:         payload.Title,
		Description:   payload.Description,
		Category:      models.AuctionCategory(payload.Category),
		City:          payload.City,
	}

	err = handler.AuctionRepo.Create(r.Context(), &auction)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create auction")
		return
	}

	response := map[string]string{
		"message":    "Auction created successfully",
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
