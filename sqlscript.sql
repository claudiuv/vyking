/*CREATE DATABASE Vyking;*/

USE Vyking;

CREATE TABLE Players (
    PlayerID INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(250) NOT NULL,
    Email VARCHAR(300) NOT NULL,
    AccountBalance DECIMAL(10, 2) NOT NULL DEFAULT 0.00
);

CREATE TABLE Tournaments (
    TournamentID INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    PrizePool DECIMAL(10, 2) NOT NULL,
    StartDate DATE NOT NULL,
    EndDate DATE NOT NULL
);

CREATE TABLE Bets (
    BetID INT AUTO_INCREMENT PRIMARY KEY,
    PlayerID INT NOT NULL,
    TournamentID INT NOT NULL,
    BetAmount DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (PlayerID) REFERENCES Players(PlayerID),
    FOREIGN KEY (TournamentID) REFERENCES Tournaments(TournamentID)
);

-- Prizes Table
CREATE TABLE Prizes (
    PrizeID INT AUTO_INCREMENT PRIMARY KEY,
    TournamentID INT NOT NULL,
    PlayerID INT NOT NULL,
    PrizeAmount DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (TournamentID) REFERENCES Tournaments(TournamentID),
    FOREIGN KEY (PlayerID) REFERENCES Players(PlayerID)
);


DELIMITER //

CREATE PROCEDURE DistributePrizes(IN tournamentId INT)
BEGIN
	IF EXISTS(SELECT * FROM Prizes WHERE Prizes.TournamentID = tournamentID) THEN
	  SET @msg = CONCAT('Prizes for the tournament with tournamentID=', tournamentID, '  are already distributed');
	  SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @msg;
	END IF;

	INSERT INTO Prizes (TournamentID, PlayerID, PrizeAmount)	
	SELECT tournamentId AS TournamentID, T.PlayerID, T.PrizeAmount
	FROM (
		SELECT
		  b.PlayerID,
		  CASE RANK() OVER (ORDER BY SUM(b.BetAmount) DESC)
		      WHEN 1 THEN t.PrizePool * 0.50
		      WHEN 2 THEN t.PrizePool * 0.30
		      WHEN 3 THEN t.PrizePool * 0.20
		      ELSE 0
		  END AS PrizeAmount	
		FROM Bets b
		INNER JOIN Tournaments t ON t.TournamentID = b.TournamentID
		WHERE b.TournamentID = tournamentId
		GROUP BY b.PlayerID, t.PrizePool
	) T	
	WHERE T.PrizeAmount > 0;

	UPDATE Players
	JOIN Prizes ON Players.PlayerID = Prizes.PlayerID
	SET Players.AccountBalance = Players.AccountBalance + Prizes.PrizeAmount
	WHERE Prizes.TournamentID = tournamentId;

END //

DELIMITER ;

